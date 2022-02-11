# flake8: noqa: E501

import os
from os import path


PKG_PATH = 'pkg'
ROOT_PKG = 'github.com/zekroTJA/shinpuru'
OUTPUT_FILE = "./docs/public-packages.md"


def get_pkg_description(pkg):
    pkg_path = path.join(PKG_PATH, pkg)
    for file in os.listdir(pkg_path):
        if not file.endswith('.go') or file.endswith('_test.go'):
            continue
        with open(path.join(pkg_path, file), 'r', encoding='utf-8') as f:
            desc = ''
            line = ''
            while not line.startswith('package '):
                desc += line[3:]
                line = f.readline()
            if not desc:
                continue
            return desc.strip()


res = \
    '<!-- insert:PUBLIC_PACKAGES -->\n' \
    '# Public Packages\n\n'
for pkg in os.listdir(PKG_PATH):
    if not path.isdir(path.join(PKG_PATH, pkg)):
        continue
    print(f'Processing package {pkg} ...')
    desc = (get_pkg_description(pkg)
            or 'No package description.').replace('\n', ' ')
    res += '- [**`{root}/{sub}/{pkg}`**]({sub}/{pkg})  \n  *{desc}*\n\n'.format_map({
        'root': ROOT_PKG,
        'sub': PKG_PATH,
        'pkg': pkg,
        'desc': desc,
    })

with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
    f.write(res)
