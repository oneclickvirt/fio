#!/usr/bin/env python3
"""Choose a stable patch version above every valid legacy native component tag."""
import datetime
import re
import subprocess


def next_version(tags):
    versions = [(0, 0, 0)]
    for tag in tags:
        match = re.fullmatch(r'v(\d+)\.(\d+)\.(\d+)(?:-(\d{14}))?', tag)
        if not match or int(match[1]) != 0:
            continue
        if match[4]:
            try:
                datetime.datetime.strptime(match[4], '%Y%m%d%H%M%S')
            except ValueError:
                continue
        versions.append(tuple(map(int, match.groups()[:3])))
    major, minor, patch = max(versions)
    return f'v{major}.{minor}.{patch + 1}'


if __name__ == '__main__':
    tags = subprocess.check_output(['git', 'tag', '--list'], text=True).splitlines()
    print(next_version(tags))
