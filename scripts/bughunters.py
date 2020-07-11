# flake8: noqa: E501

import requests
import json
import os
import codecs
from typing import Dict, List


GH_GQL_ENDPOINT = 'https://api.github.com/graphql'
OUTPUT_FILE = './bughunters.md'

HEADERS = {
    'Authorization': "bearer " + os.environ.get('GITHUB_TOKEN')
}


def do_req(query: str) -> Dict:
    data = json.dumps({
        'query': query
    })
    r = requests.post(GH_GQL_ENDPOINT, data, headers=HEADERS)
    return r.json()


def query_issues(login: str, repo: str, after: str = None) -> Dict:
    after = ', after: {}'.format(after) if after else ''
    query = '''
    query {
      user(login: "%s") {
        repository(name: "%s") {
          issues(first: 100%s) {
            totalCount,
            edges {
              node {
                author {
                  login
                },
                number,
                url
              }
            }
          }
        }
      }
    }
    '''
    res = do_req(query % (login, repo, after))
    return res


def query_prs(login: str, repo: str, after: str = None) -> Dict:
    after = ', after: {}'.format(after) if after else ''
    query = '''
    query {
      user(login: "%s") {
        repository(name: "%s") {
          pullRequests(first: 100%s) {
            totalCount,
            edges {
              node {
                author {
                  login
                },
                number,
                url
              }
            }
          }
        }
      }
    }
    '''
    res = do_req(query % (login, repo, after))
    return res


def query_all_issues(login: str, repo: str) -> List[Dict]:
    issues = []
    after = None
    while True:
        res = query_issues(login, repo, after)
        iss_res = res.get("data").get("user").get("repository").get("issues")
        n = iss_res.get("totalCount")
        issues += [e.get("node") for e in iss_res.get("edges")]
        if n < 100:
            break

    return issues


def query_all_prs(login: str, repo: str) -> List[Dict]:
    prs = []
    after = None
    while True:
        res = query_prs(login, repo, after)
        prs_res = res.get("data").get("user").get("repository").get("pullRequests")
        n = prs_res.get("totalCount")
        prs += [e.get("node") for e in prs_res.get("edges")]
        if n < 100:
            break

    return prs


if __name__ == "__main__":
    issues = query_all_issues("zekroTJA", "shinpuru")
    issues = [i for i in issues if i.get("author") and i.get("author").get("login") != "zekroTJA"]
    
    prs = query_all_prs("zekroTJA", "shinpuru")
    prs = [p for p in prs if p.get("author") and p.get("author").get("login") != "zekroTJA"]

    bhs = {}

    for i in issues:
        author = i.get("author").get("login")
        if author not in bhs:
            bhs[author] = []
        bhs[author].append(i)

    for p in prs:
        author = p.get("author").get("login")
        if author not in bhs:
            bhs[author] = []
        bhs[author].append(p)

    data = "# Bug Hunters\n\n" + \
           "A list to honor all people who found some bugs, had some great ideas or contributed directly to shinpuru. ❤️\n\n" + \
           "| GitHub | Issues / PRs |\n" + \
           "|--------|--------------|\n"

    items = sorted(bhs.items(), key=lambda x: x[0].lower())
    for k, v in items:
        v.sort(key=lambda n: int(n.get("number")))
        nums = ["[#{}]({})".format(n.get("number"), n.get("url")) for n in v]
        data += "| [{}](https://github.com/{}) | {} |\n".format(k, k, ', '.join(nums))

    with codecs.open(OUTPUT_FILE, 'w', 'utf-8') as f:
        f.write(data)