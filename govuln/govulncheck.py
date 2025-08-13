import argparse
import dataclasses
import json
import os
import subprocess
import sys
from typing import Tuple

from internals import Internal


class Table:
    def __init__(self, title, *headers):
        self.title = title
        self.headers = headers
        self.columns = [len(h) for h in headers]
        self.rows = []

    def write(self, *data):
        for i, d in enumerate(data):
            self.columns[i] = max(self.columns[i], len(d))
        self.rows.append(data)

    def _row(self, *data):
        print('| ', end='')
        for i, d in enumerate(data):
            print(f'{d:<{self.columns[i]}}', end=' | ')
        print()

    def _line(self):
        data = ['-' * d for d in self.columns]
        self._row(*data)

    def print(self):
        print(f'| {self.title:<{sum(self.columns) + len(self.columns) * 2}} |')
        self._line()
        self._row(*self.headers)
        self._line()
        for row in self.rows:
            self._row(*row)


@dataclasses.dataclass
class Affected:
    package: str
    introduced: str
    fixed: str

    def __str__(self) -> str:
        return f'{self.package} Introduced:{self.introduced} Fixed:{self.fixed}'


class Version:
    def __init__(self, v: str):
        try:
            v = str(v)
            if v.startswith('v'):
                v = v[1:]
            elif v.startswith('go'):
                v = v[2:]
            parts = v.split('-', 1)
            self.extra = parts[1] if len(parts) > 1 else ''
            values = parts[0].split('.')
            self.major = int(values[0]) if values[0].isnumeric() else values[0]
            self.minor = int(values[1]) if values[1].isnumeric() else values[1]
            if len(values) > 2:
                self.patch = int(values[2]) if values[2].isnumeric() else values[2]
            else:
                self.patch = 0
        except Exception:
            self.major = 0
            self.minor = 0
            self.patch = 0

    def __lt__(self, other):
        if not isinstance(other, Version):
            return NotImplemented
        return (self.major, self.minor, self.patch, self.extra) < (
            other.major,
            other.minor,
            other.patch,
            self.extra,
        )

    def __repr__(self):
        return f'{self.major}.{self.minor}.{self.patch}' + (
            '' if not self.extra else f'-{self.extra}'
        )


class OSV:
    def __init__(self, value: dict):
        self.id = value.get('id')
        self.summary = value.get('summary')
        self.affected = dict()

        for aff in value.get('affected', []):
            fixed = None
            introduced = None
            if not (package := aff.get('package', {}).get('name', '')):
                continue
            if package == 'stdlib':
                package = 'GO stdlib'
            for ranges in aff.get('ranges', []):
                for event in ranges.get('events', []):
                    if 'fixed' in event:
                        fixed = event['fixed']
                    if 'introduced' in event:
                        introduced = event['introduced']
                if fixed and introduced:
                    break
            self.affected[package] = Affected(package, introduced, fixed)

    def __repr__(self):
        s = f'{self.id}\n\t{self.summary}'
        for a in self.affected.values():
            s += f'\n\t{a}'
        return s


class SBOM:
    def __init__(self, value: dict, internal_owners: list[str]):
        self.go_version = Version(value.get('go_version', '0.0.0'))
        self.package = None
        self.modules = {'GO stdlib': self.go_version}
        self.internals: list[Internal] = []
        for module in value.get('modules', []):
            if path := module.get('path'):
                if version := module.get('version'):
                    self.modules[path] = Version(version)
                    internal = Internal(path, version)
                    if internal.valid and internal.owner in internal_owners:
                        internal.run()
                        if internal.has_vulnerabilities:
                            self.internals.append(internal)

                elif not self.package:
                    self.package = path


class GoVulnCheck:
    def __init__(
        self, just_warn: bool = False, internal_repositories_owners: list[str] = []
    ):
        self.sbom: SBOM = None
        self.osvs: list[OSV] = []
        self.just_warn = just_warn
        self.internal_repositories_owners = internal_repositories_owners

    def summary(self):
        if not self.osvs:
            return

    def call(self) -> bool:
        if not (self.check_dependencies() and self.run_vulncheck()):
            return False
        if self.sbom.internals:
            table = Table(
                'Affected internal repositories',
                'Version',
                'Package',
                'Vulnerabilities',
            )
            for s in self.sbom.internals:
                vulns = ', '.join([v['code'] for v in s.meta.vulnerabilities])

                table.write(str(Version(s.version)), s.repo, vulns)

            table.print()
            print()

        if not self.osvs:
            return True

        affected = dict()
        for osv in self.osvs:
            for af in osv.affected.values():
                fixed = Version(af.fixed)
                package = affected.setdefault(
                    af.package, dict(package=af.package, fixed=fixed, vulns=[])
                )
                package['vulns'].append(osv)
                if package['fixed'] < fixed:
                    package['fixed'] = fixed

        table = Table('Affected packages', 'Current', 'Fixed', 'Package')
        for package, data in affected.items():
            if not (cur_version := self.sbom.modules.get(package)):
                continue
            if cur_version < data['fixed']:
                table.write(str(cur_version), str(data['fixed']), package)

        table.print()

        print('\nℹ️  Run "govulncheck ./..." for more details')
        return False

    def __call__(self) -> int:
        if self.call():
            return 0
        if self.just_warn:
            print('\n⚠️  Just a warning. The commit will not be blocked.')
            return 0
        return 1

    def _run_command(self, *command) -> Tuple[int, str]:
        result = subprocess.run(command, capture_output=True, check=False)
        output = (
            result.stderr.decode('utf-8') if result.returncode else ''
        ) + result.stdout.decode('utf-8')
        if result.returncode != 0:
            print(f'Running command "{command}" resulted {result.returncode}')
        return result.returncode, output

    def check_dependencies(self) -> bool:
        """Checks if go.mod exists and govulncheck is available"""
        if not os.path.isfile('go.mod'):
            print('missing go.mod file')
            return False
        # Check govulncheck
        rc, _ = self._run_command('govulncheck', '-json')
        if rc == 0:
            return True

        # Try install govulncheck
        cmd = ('go', 'install', 'golang.org/x/vuln/cmd/govulncheck@latest')
        rc, output = self._run_command(*cmd)
        if rc == 0:
            return True
        print('Error installing govulncheck: ', ' '.join(cmd))
        print(output)
        return False

    def run_vulncheck(self) -> bool:
        """Runs govulncheck and parses data. Returns false if there are any execution error"""
        rc, output = self._run_command('govulncheck', '-json', './...')
        if rc:
            print('Error running govulncheck: ', rc, output)
            return False
        body = ''
        for line in output.splitlines():
            if line.startswith('{'):
                body = line
            elif line.startswith('}'):
                body += line
                if error := self.parse(body):
                    print('Error parsing govulncheck: ', error)
                    return False
                body = ''
            else:
                body += line
        return True

    def parse(self, body) -> str:
        try:
            field = json.loads(body)
            if not (isinstance(field, dict) and len(field) == 1):
                return ''
            key = list(field.keys())[0]
            if key == 'osv':
                self.osvs.append(OSV(field[key]))
            elif key == 'SBOM':
                self.sbom = SBOM(field[key], self.internal_repositories_owners)
            return ''
        except Exception as exc:
            return str(exc)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--just-warn',
        action='store_true',
        default=False,
        help='Doesn`t block the commit, just warn',
    )
    parser.add_argument(
        '--internal-repositories-owners',
        nargs='+',
        default=['melisource', 'mercadolibre'],
        help='List of internal repositories owners',
    )

    args = parser.parse_args()
    govuln = GoVulnCheck(args.just_warn, args.internal_repositories_owners)
    return govuln()


if __name__ == '__main__':
    sys.exit(main())
