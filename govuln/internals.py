import json
import os
import pathlib
import subprocess
from datetime import datetime, timedelta


class Internal:
    root = pathlib.Path().home() / '.govuln'

    CLONED = 'cloned'
    CHECKED = 'checked'
    CLONE_ERROR = 'clone_error'

    class Meta:
        DATE_FORMAT = '%Y-%m-%d %H:%M:%S'

        def __init__(self, filename):
            self.filename = filename
            self.last_update: datetime = datetime(1, 1, 1)
            self.status: str = ''
            self.vulnerabilities: list[dict] = []
            if not self.load():
                self.save()

        def load(self) -> bool:
            if os.path.isfile(self.filename):
                try:
                    with open(self.filename) as f:
                        content = json.load(f)
                    if not isinstance(content, dict):
                        raise ValueError('Invalid content for file', self.filename)
                    if not (last_update := content.get('last_update')):
                        raise ValueError('Missing field last_update', self.filename)
                    self.last_update = datetime.strptime(last_update, self.DATE_FORMAT)
                    if not (status := content.get('status')):
                        raise ValueError('Missing field status', self.filename)
                    self.status = status
                    self.vulnerabilities = content.get('vulnerabilities', [])
                    return True
                except Exception:
                    ...

        def save(self):
            try:
                with open(self.filename, 'w') as f:
                    json.dump(
                        dict(
                            last_update=self.last_update.strftime(self.DATE_FORMAT),
                            status=self.status,
                            vulnerabilities=self.vulnerabilities,
                        ),
                        f,
                    )
            except Exception:
                ...

    def __init__(self, path: str, version: str):
        """path:  {
          "path": "github.com/mattn/go-isatty",
          "version": "v0.0.20"
        },
        {
          "path": "github.com/melisource/fury_fbm-fiscal-coverage-toolkit/v2",
          "version": "v2.6.0"
        },"""
        self.valid = False
        parts = path.split('/')
        if len(parts) < 3:
            # Should be a repository with host, owner and name
            return

        domain, owner, repository = parts[0:3]
        self.owner = owner
        self.repo = f'{domain}/{owner}/{repository}'
        self.extra_path = '/'.join(parts[3:])
        self.version = version
        self.path = self.root / domain / owner / repository / version
        try:
            os.makedirs(self.path, exist_ok=True)

        except Exception as exc:
            print('Failed to create folder: ', self.path, exc)
            return
        self.meta = self.Meta(
            self.root / domain / owner / repository / f'{version}.json'
        )
        self.valid = True
        self.cloned = (
            os.path.isdir(self.path / '.git') or self.meta.status == self.CLONE_ERROR
        )
        self.giturl = f'git@{domain}:{owner}/{repository}.git'
        self.has_vulnerabilities = bool(self.meta.vulnerabilities)

    def _run(self, *command) -> tuple[int, str]:
        result = subprocess.run(command, capture_output=True, check=False)
        output = result.stderr.decode('utf-8') + result.stdout.decode('utf-8')
        return result.returncode, output

    def clone(self) -> bool:
        # git clone --branch v2.6.0 git@github.com:melisource/fury_fbm-fiscal-coverage-toolkit.git fc
        if self.cloned:
            return True
        print('Cloning ', self.giturl, '...')
        result, output = self._run(
            'git', 'clone', '--branch', self.version, self.giturl, str(self.path)
        )
        if result == 0:
            # Cloned with success
            return True
        else:
            if 'Remote branch' in output and 'not found in upstream origin' in output:
                self.meta.status = self.CLONE_ERROR
                self.meta.save()
            print('Error cloning', self.giturl)

    def check(self):
        print('Checking vulnerabilities', self.giturl, '...', end='')
        result, output = self._run('govulncheck', '-C', str(self.path), './...')
        vulnerabilities = []
        if result:
            code = ''
            description = ''
            fixed_in = ''
            for line in output.splitlines():
                if line.startswith('Vulnerability #'):
                    code = line.split(': ')[1].strip()
                elif code and not description:
                    description = line.strip()
                elif line.startswith('    Fixed in:'):
                    fixed_in = line.split(': ')[1].strip()
                    vulnerabilities.append(
                        dict(code=code, description=description, fixed_in=fixed_in)
                    )
                    code = ''
                    description = ''
            print('found', len(vulnerabilities), 'vulnerabilities')
        else:
            print('no vulnerabilities found')

        self.meta.last_update = datetime.now()
        self.meta.vulnerabilities = vulnerabilities
        self.meta.status = self.CHECKED
        self.meta.save()

    def run(self):
        if self.has_vulnerabilities:
            return True
        if not self.cloned:
            if self.clone():
                self.meta.last_update = datetime.now()
                self.meta.status = self.CLONED
                self.meta.save()
        if self.meta.last_update < datetime.now() - timedelta(days=1):
            return True
        if self.check():
            self.meta.last_update = datetime.now()
            self.meta.status = self.CHECKED
            self.meta.save()
            return True


if __name__ == '__main__':
    i = Internal('github.com/melisource/fury_fbm-fiscal-coverage-toolkit/v2', 'v2.6.0')
