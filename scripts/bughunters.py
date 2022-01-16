# flake8: noqa: E501

import itertools
import requests
import json
import os
import codecs
from typing import Dict, List


GH_GQL_ENDPOINT = 'https://api.github.com/graphql'
OUTPUT_FILE = './bughunters.md'

HEADERS = {
    'Authorization': f"bearer {os.environ.get('GITHUB_TOKEN')}"
}

POINTS_FOR_ISSUE = 1
POINTS_FOR_PR = 2


class BHEntry:
    def __init__(self):
        self.issues = []
        self.prs = []
        self.points = 0

    def add_issue(self, issue):
        self.issues.append(issue)
        self.points += POINTS_FOR_ISSUE

    def add_pr(self, pr):
        self.prs.append(pr)
        self.points += POINTS_FOR_PR

    def sort_items(self):
        self.issues.sort(key=lambda n: int(n.get("number")))
        self.prs.sort(key=lambda n: int(n.get("number")))

    def get_formatted_issues(self):
        return [f'[#{n.get("number")}]({n.get("url")})' for n in self.issues]

    def get_formatted_prs(self):
        return [f'[#{n.get("number")}]({n.get("url")})' for n in self.prs]


def do_req(query: str) -> Dict:
    data = json.dumps({
        'query': query
    })
    r = requests.post(GH_GQL_ENDPOINT, data, headers=HEADERS)
    return r.json()


def query_issues(login: str, repo: str, after: str = None) -> Dict:
    after = f', after: "{after}"' if after else ''
    query = '''
    query {
      user(login: "%s") {
        repository(name: "%s") {
          issues(first: 100%s) {
            totalCount,
            edges {
              cursor,
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
    after = f', after: "{after}"' if after else ''
    query = '''
    query {
      user(login: "%s") {
        repository(name: "%s") {
          pullRequests(first: 100%s) {
            totalCount,
            edges {
              cursor,
              node {
                author {
                  login
                },
                number,
                url,
                merged,
                changedFiles,
                additions,
                deletions
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
        edges = iss_res.get("edges")
        n = len(edges)
        issues += [e.get("node") for e in edges]
        if n < 100:
            break
        after = edges[-1].get("cursor")

    return issues


def query_all_prs(login: str, repo: str) -> List[Dict]:
    prs = []
    after = None
    while True:
        res = query_prs(login, repo, after)
        prs_res = res.get("data").get("user").get(
            "repository").get("pullRequests")
        edges = prs_res.get("edges")
        n = len(edges)
        prs += [e.get("node") for e in edges]
        if n < 100:
            break
        after = edges[-1].get("cursor")

    return prs


def medal(i: int) -> str:
    if i > 2:
        return ""
    return ["🥇", "🥈", "🥉"][i] + " "


def get_pr_stats(prs: List[Dict]) -> Dict:
    additions = 0
    deletions = 0
    changedFiles = 0
    for pr in prs:
        additions += pr.get('additions')
        deletions += pr.get('deletions')
        changedFiles += pr.get('changedFiles')
    return [additions, deletions, changedFiles]


if __name__ == "__main__":
    issues = query_all_issues("zekroTJA", "shinpuru")
    issues = [i for i in issues if i.get("author") and i.get(
        "author").get("login") != "zekroTJA"]

    prs = query_all_prs("zekroTJA", "shinpuru")
    prs = [p for p in prs if p.get("merged") and p.get(
        "author") and p.get("author").get("login") != "zekroTJA"]

    bhs = {}

    for i in issues:
        author = i.get("author").get("login")
        if author not in bhs:
            bhs[author] = BHEntry()
        bhs[author].add_issue(i)

    for p in prs:
        author = p.get("author").get("login")
        if author not in bhs:
            bhs[author] = BHEntry()
        bhs[author].add_pr(p)

    [additions, deletions, changedFiles] = get_pr_stats(prs)

    data = "# Bug Hunters\n\n" \
           "A list to honor all people who found some bugs, had some great ideas " \
           "or contributed directly to shinpuru. ❤️\n\n" \
           f"In total, **{len(bhs)}** different wonderful people contributed a sum of " \
           f"**{len(issues)}** issues and **{len(prs)}** pull requests (with {additions} " \
           f"added and {deletions} deleted lines of code in {changedFiles} different files)! 🎉\n\n" \
           "| GitHub | Issues | PRs | Points* |\n" \
           "|--------|--------|-----|---------|\n"

    items = sorted(bhs.items(), key=lambda x: x[1].points, reverse=True)

    i = 0
    for k, v in items:
        v.sort_items()
        issues = v.get_formatted_issues()
        prs = v.get_formatted_prs()
        points = v.points
        data += f'| {medal(i)} [{k}](https://github.com/{k}) | {", ".join(issues)} | {", ".join(prs)} | `{points}` |\n'
        i += 1

    data += f'\n\n---\n*For explaination: A contributor gets `{POINTS_FOR_ISSUE}` point(s) for each ' \
        'created issue and `{POINTS_FOR_PR}` point(s) for each **merged** pull request.'

    with codecs.open(OUTPUT_FILE, 'w', 'utf-8') as f:
        f.write(data)
