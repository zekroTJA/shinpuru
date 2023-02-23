import argparse
import os
from os import path
import json


def parse_args():
    p = argparse.ArgumentParser()

    p.add_argument('-b', '--base', type=str, required=True,
                   help="The location of the base language files")
    p.add_argument('-t', '--target', type=str, required=True,
                   help="Target language directories")

    return p.parse_args()


def main():
    args = parse_args()

    for (root, _, files) in os.walk(args.base):
        for file in files:
            file = path.join(root, file)
            base_data = {}
            target_data = {}

            with open(file) as f:
                base_data = json.load(f)

            target_file = file.replace(args.base, args.target)
            if path.exists(target_file):
                with open(target_file) as f:
                    target_data = json.load(f)

            merge(base_data, target_data)

            os.makedirs(path.dirname(target_file), exist_ok=True)
            with open(target_file, mode='w') as f:
                json.dump(target_data, f, indent=2, ensure_ascii=False)


def merge(base: dict, target: dict):
    for (k, v) in base.items():
        if type(v) == dict:
            if k not in target:
                target[k] = {}
            merge(v, target[k])
        else:
            if k not in target or target[k] == "":
                target[k] = v


if __name__ == '__main__':
    main()
