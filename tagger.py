import io
import sys
import json

from urllib.request import urlopen


#_URL = "http://en.wiktionary.org/wiki/"
_URL = "http://en.wiktionary.org/w/api.php"
URL = _URL + "?action=parse&page=%s&prop=sections|categories|images&format=json"


class Dictionary:

    def parts_of_speech(self, word):

        f = urlopen(URL % word)
        response = str(f.read(), encoding="utf-8")
        jo = json.loads(response)
        sections = jo['parse']['sections']
        seen = set()
        for section in sections:
            line = section['line']
            if line in ["Noun", "Adjective", "Verb"]:
                if line not in seen:
                    seen.add(line)
                    yield line

DICTIONARY = Dictionary()


if __name__ == "__main__":
    d = Dictionary()
    #print(list(d.parts_of_speech("chili")))
    #print(list(d.parts_of_speech("celeriac")))
    print("> ", end="")
    sys.stdout.flush()
    line = sys.stdin.readline().strip()
    print(list(d.parts_of_speech(line)))
