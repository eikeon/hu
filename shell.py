import io
import sys

from chart_parser import Chart, PartOfSpeech
from tokenizer import tokenize, Token
from tokenizer import *


if __name__=="__main__":

    parts_of_speech = [
        PartOfSpeech("SENTENCE", [WORD, WORD, STRING, PERIOD, NEWLINE]),
        PartOfSpeech("SENTENCE", [WORD, WORD, WORD, PERIOD, NEWLINE]),
        PartOfSpeech("SENTENCE", ["NOUN", "VERB", "NOUN", PERIOD, NEWLINE]),
        PartOfSpeech("SENTENCE", [Token("I", WORD), WORD, STRING, PERIOD, NEWLINE]),
        PartOfSpeech("SENTENCE", ["SUBJECTIVE_NOUN", "VERB", "OBJECT", PERIOD, NEWLINE]),
        PartOfSpeech("OBJECT", ["NOUN"]),
        PartOfSpeech("OBJECT", ["ADJECTIVE", "NOUN"])
        #PartOfSpeech("UNKNOWN", [Token]),
        #PartOfSpeech("UNKNOWN", [Token, "UNKNOWN"]),
        ]
    
    c = Chart(rules=parts_of_speech)

    while True:
        print("> ", end="")
        sys.stdout.flush()
        line = io.StringIO(sys.stdin.readline())
        tokens = list(tokenize(line.__next__))


        print("Tokens:")
        for token in tokens:
            #print("    %r" % token)
            print("    %s(%s)" % (token.token_type, token.string))
        print("")

        print("Results:")
        for result in c.parse(tokens):
            print("    ", result)
        print("")
        print(c.num)
