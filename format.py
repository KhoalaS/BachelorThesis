import re
import os

f = open("notes.MD")
text = ""

for line in f:
    text += line

r_0 = r"(\$\$)\s(\\begin\{align\*\}\s[^$]*)(\$\$)"
r_1 = re.compile(r"(\$\$)\s(\\begin\{algorithm\}\[H\].*?end{algorithm})\s(\$\$)", re.DOTALL)
r_2 = re.compile(r"(\$\$)\s(\\begin\{algorithm\}.*?end{algorithm})\s(\$\$)", re.DOTALL)

out = open("./out/out.MD", "w+")
text = re.sub(r_0, r"\2", text)
text = re.sub(r_1, r"\2", text)
text = re.sub(r_2, r"\2", text)

out.write(text)
out.close()

code = os.system(
    "pandoc --verbose -H header.tex --number-sections --toc -V geometry:margin=1in out/out.MD -o notes.pdf --bibliography=./lit.bib --csl=./ieee.csl")
exit(code)
