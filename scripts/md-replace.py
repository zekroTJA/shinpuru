import argparse
import re


INSERT_RX = r'^<!--\s?insert:(\w+)\s?-->'
START_RX = r'<!--\s?start:(\w+)\s?-->'
END_RX = r'<!--\s?end:(\w+)\s?-->'


def parse_apply(name):
    with open(name, 'r', encoding='utf-8') as f:
        data = f.read()
        match = re.match(INSERT_RX, data)
        if not match:
            raise Exception(
                'Could not find insert marker at '
                f'the beginning of the file "{name}".')
        marker = re.findall(INSERT_RX, data)[0]
        data = data[match.span()[1]:]
        return (marker, data)


def main():
    parser = argparse.ArgumentParser('md-replace')
    parser.add_argument('--input', '-i', type=str, required=True,
                        help='The input file to apply insertions on.')
    parser.add_argument('apply', type=str, nargs='+',
                        help='The fiels to apply on the input.')
    parser.add_argument('--dry', '-d', action='store_const', const=True,
                        default=False, help='Only print the output and do '
                        'not write back to input.')
    args = parser.parse_args()

    applies = {}

    for apply_name in args.apply:
        (marker, data) = parse_apply(apply_name)
        applies[marker] = data

    input_data = ''
    with open(args.input, 'r', encoding='utf-8') as f:
        input_data = f.read()

    blocks = {}

    for s in re.finditer(START_RX, input_data):
        (start, end) = s.span()
        marker = re.findall(START_RX, input_data[start:end])[0]
        blocks[marker] = [(start, end), ()]

    for s in re.finditer(END_RX, input_data):
        (start, end) = s.span()
        marker = re.findall(END_RX, input_data[start:end])[0]
        block = blocks.get(marker)
        if not block:
            continue
        if start < block[0][0]:
            raise Exception(
                f'Block "{marker}" has the end '
                'marker before the start marker.')
        blocks[marker][1] = (start, end)

    offset = 0
    for marker, block in blocks.items():
        apply = applies.get(marker)
        if not apply:
            continue

        start = block[0][1]
        end = block[1][0]
        input_data = input_data[:start+offset] + \
            apply + input_data[end+offset:]

        offset += len(apply) - (end - start)

    if args.dry:
        print(input_data)
    else:
        with open(args.input, 'w', encoding='utf-8') as f:
            f.write(input_data)


if __name__ == '__main__':
    main()
