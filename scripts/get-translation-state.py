#!/usr/bin/env python3

import argparse
import os
from os import path
import json
import re


WORDS_RX = r'[\w]+'


def parse_args():
    p = argparse.ArgumentParser()

    p.add_argument('dir', nargs=1, type=str,
                   help="Language pack directory.")
    p.add_argument('--base', '-b', type=str, default="en-US",
                   help="Base language code to compare against.")

    return p.parse_args()


def main() -> int:
    args = parse_args()

    indices = {}
    for dir in os.listdir(args.dir[0]):
        indices[dir] = index_trans_pack(path.join(args.dir[0], dir))

    base = indices.get(args.base)
    if not base:
        print("Error: Given base translation pack does not exist")
        return 1

    stats = {}
    for (key, index) in indices.items():
        files = len(index)
        (keys, words) = get_n_kvs(index)
        stats[key] = (files, keys, words)

    base_state = stats[args.base]
    for (key, stats) in [(k, s) for (k, s) in stats.items() if k != args.base]:
        print_state(key, base_state, stats)

    return 0


def index_trans_pack(dir):
    index = {}
    for (root, _, files) in os.walk(dir):
        for file in files:
            with open(path.join(root, file)) as f:
                index[file] = json.load(f)
    return index


def get_n_kvs(index):
    n_keys = 0
    n_words = 0
    for (_, v) in index.items():
        if type(v) == dict:
            (keys, words) = get_n_kvs(v)
            n_keys += keys
            n_words += words
        elif type(v) == list or type(v) == tuple:
            n_keys += len(v)
            n_words += sum([count_words(e) for e in v])
        elif len(v) != 0:
            n_keys += 1
            n_words += count_words(v)
    return (n_keys, n_words)


def count_words(s):
    return len(re.findall(WORDS_RX, s))


def print_state(key, base, target):
    (b_files, b_keys, b_words) = base
    (t_files, t_keys, t_words) = target
    print(
        f"{key}:\n"
        f"  files: {t_files:>5} / {b_files:>5} ({t_files/b_files:.1%})\n"
        f"  keys:  {t_keys:>5} / {b_keys:>5} ({t_keys/b_keys:.1%})\n"
        f"  words: {t_words:>5} / {b_words:>5} ({t_words/b_words:.1%})\n"
    )


if __name__ == "__main__":
    exit(main())
