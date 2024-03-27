import re
import os

f = open("./tex/notes.MD")
text = ""

for line in f:
    text += line

r = re.compile(r"(\$\$)\s(\\begin.*?)(\$\$)", re.DOTALL)

out = open("./out/out.MD", "w+")
text = re.sub(r, r"\2", text)

out.write(text)
out.close()

code = os.system(
    "pandoc --verbose -H ./tex/header.tex --template ./tex/template.tex --number-sections --toc -V geometry:margin=1in out/out.MD -o notes.pdf --bibliography=./tex/lit.bib --csl=./tex/ieee.csl")
exit(code)
