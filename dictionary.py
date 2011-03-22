from urllib.request import urlopen


_URL = "http://www.google.com/dictionary?hl=en&sl=en&tl=en&q="
# http://dictionary.reference.com/api/

class Dictionary:

    def part_of_speech(self, word):
        f = urlopen(_URL+word)
        s = f.read()
        return s

if __name__ == "__main__":
    d = Dictionary()
    print(d.part_of_speech("black"))
