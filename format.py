import re
import os

f = open("README.MD")
text = "# This document was auto-generated using pandoc.\n\n"

header = (
    "---",
    "title: Implementation and evaluation of a self-monitoring approximation algorithm for 3-Hitting-Set",
    "header-includes: |",
    "\t\\usepackage{algorithm2e}",
    "...")

latex_commands = ["\RestyleAlgo{ruled}", "\SetAlgoLined",
                  "\DontPrintSemicolon"]

text = text + "\n".join(header) + "\n"
text = text + "\n".join(latex_commands) + "\n"


for line in f:
    text += line

r_0 = r"(\$\$)\s(\\begin\{align\*\}\s[^$]*)(\$\$)"
r_1 = re.compile(r"(\$\$)\s(\\begin\{algorithm\}\[H\].*)(\$\$)", re.DOTALL)

out = open("./out/out.MD", "w+")
new = re.sub(r_0, r"\2", text)
new = re.sub(r_1, r"\2", new)
out.write(new)


out.close()

code = os.system(
    "pandoc -V geometry:margin=1in out/out.MD -o notes.pdf --bibliography=./lit.bib --csl=./ieee.csl")
exit(code)
