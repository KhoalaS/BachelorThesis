import re
import os

f = open("README.MD")
text = "# This document was auto-generated using pandoc.\n\n"
for line in f:
    text += line

r = r"(\$\$)\s(\\begin\{align\*\}\s[^$]*)(\$\$)"

new = re.sub(r, r"\2", text)

out = open("./out/out.MD", "w+")
out.write(new)
out.close()

code = os.system("pandoc -V geometry:margin=1in out/out.MD -o notes.pdf --bibliography=./lit.bib --csl=./ieee.csl")
exit(code)
