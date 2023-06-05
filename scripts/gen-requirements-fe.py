import json
import requests


OUTPUT = './docs/requirements-fe.md'


def main():
    data = {}
    with open('web/package.json') as f:
        data = json.load(f)

    modules = []
    for (name, version) in data.get("dependencies").items():
        print(f"Processing package {name} ...")
        resp = requests.get(f'https://registry.npmjs.com/{name}/latest')
        homepage = resp.json().get('homepage') or \
            f"https://www.npmjs.com/package/{name}"
        modules.append(f"[{name}]({homepage}) `({version})`")

    with open(OUTPUT, 'w') as f:
        f.write('<!-- insert:REQUIREMENTS_FE -->\n')
        f.writelines([f'- {m}\n' for m in modules])


if __name__ == '__main__':
    main()
