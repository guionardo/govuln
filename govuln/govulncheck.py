import argparse
import dataclasses
import json
import os
import subprocess
from typing import Generator

from packaging.version import Version
import sys


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


# class Version:
#     def __init__(self, v: str):
#         try:
#             v = str(v)
#             if v.startswith('v'):
#                 v = v[1:]
#             elif v.startswith('go'):
#                 v = v[2:]
#             parts = v.split('-', 1)
#             self.extra = parts[1] if len(parts) > 1 else ''
#             values = parts[0].split('.')
#             self.major = int(values[0]) if values[0].isnumeric() else values[0]
#             self.minor = int(values[1]) if values[1].isnumeric() else values[1]
#             if len(values) > 2:
#                 self.patch = int(values[2]) if values[2].isnumeric() else values[2]
#             else:
#                 self.patch = 0
#         except Exception:
#             self.major = 0
#             self.minor = 0
#             self.patch = 0

#     def __lt__(self, other):
#         if not isinstance(other, Version):
#             return NotImplemented
#         return (self.major, self.minor, self.patch, self.extra) < (
#             other.major,
#             other.minor,
#             other.patch,
#             self.extra,
#         )

#     def __repr__(self):
#         return f'{self.major}.{self.minor}.{self.patch}' + (
#             '' if not self.extra else f'-{self.extra}'
#         )


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
    def __init__(self, value: dict):
        self.go_version = Version(value.get('go_version', '0.0.0'))
        self.package = None
        self.modules = {'GO stdlib': self.go_version}
        for module in value.get('modules', []):
            if path := module.get('path'):
                if version := module.get('version'):
                    self.modules[path] = Version(version)
                elif not self.package:
                    self.package = path


def _run_command(*command):
    result = subprocess.run(command, capture_output=True, check=False)
    output = (str(result.stderr) if result.returncode else '') + str(result.stdout)
    if result.returncode != 0:
        print(f'Running command "{command}" resulted {result.returncode}')

    return result.returncode, output


def _check_dependencies() -> bool:
    if not os.path.isfile('go.mod'):
        print('missing go.mod file')
        return False
    # Check govulncheck
    rc, output = _run_command('govulncheck', '-json')
    if rc == 0:
        return True

    # Try install govulncheck
    cmd = ('go', 'install', 'golang.org/x/vuln/cmd/govulncheck@latest')
    rc, output = _run_command(*cmd)
    if rc == 0:
        return True
    print('Error installing govulncheck: ', ' '.join(cmd))
    print(output)


def _run_vulncheck() -> Generator[str, None, None]:
    cmd = ['govulncheck', '-json', './...']
    print('Running govulncheck')
    popen = subprocess.Popen(cmd, stdout=subprocess.PIPE, universal_newlines=True)
    for stdout_line in iter(popen.stdout.readline, ''):
        yield stdout_line
    popen.stdout.close()
    return_code = popen.wait()
    if return_code:
        raise subprocess.CalledProcessError(return_code, cmd)


def parse_key(key, value, osvs: list[OSV]):
    global sbom
    if key == 'osv':
        osvs.append(OSV(value))
    elif key == 'SBOM':
        sbom = SBOM(value)


def parse(body, osvs):
    try:
        field = json.loads(body)
        if isinstance(field, dict) and len(field) == 1:
            key = list(field.keys())[0]
            value = field[key]
            parse_key(key, value, osvs)

    except Exception:
        pass  # TODO: Tratar erro no parsing


def run(osvs):
    body = ''
    for line in _run_vulncheck():
        if line.startswith('{'):
            body = line
        elif line.startswith('}'):
            body += line
            parse(body, osvs)
            body = ''
        else:
            body += line


sbom: SBOM = None


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--just-warn',
        action='store_true',
        default=False,
        help='Doesn`t block the commit, just warn',
    )

    args = parser.parse_args()

    if not _check_dependencies():
        return 1

    osvs: list[OSV] = []
    run(osvs)
    affected = dict()
    if not osvs:
        return 0

    for osv in osvs:
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
        if not (cur_version := sbom.modules.get(package)):
            continue
        if cur_version < data['fixed']:
            table.write(str(cur_version), str(data['fixed']), package)

    table.print()
    print('\nRun "govulncheck -json ./..." for more details')

    if args.just_warn:
        print('\nJust warning, not blocking the commit')
        return 0
    return 1


if __name__ == '__main__':
    sys.exit(main())
