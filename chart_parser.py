from tokenizer import Token


class PartOfSpeech(object):
    """A PartOfSpeech represents a phrase-structure production rule of the form:
            label => neighbors
       where label expands to neighbors (and therefore neighbors can
       be compsed into a label.
    """
    def __init__(self, label, neighbors=None):
        self.label = label
        self.neighbors = neighbors

    def __repr__(self):
        if self.neighbors:
            return '%s: %s' % (self.label, " ".join(self.neighbors))
        else:
            return self.label

    def matches(self, category):
        return self.neighbors and self.neighbors[0] == category


class Constituent(object):
    """An element considered as part of a construction."""

    def __init__(self, label, children, left, right):
        self.label = label
        self.children = children
        self.left = left
        self.right = right

    def __repr__(self):
        return '%s<%s>' % (self.label, repr(self.children))

    def tokens(self):
        for c in self.children:
            for t in c.tokens():
                yield t

    def terminals(self):
        for c in self.children:
            for terminal in c.terminals():
                yield terminal


class PreTerminal(Constituent):
    def __init__(self, label, token, left):
        Constituent.__init__(self, label, None, left, left+1)
        self.token = token
        
    def __repr__(self):
        return '%s(%s)' % (self.label, self.token)
        
    def tokens(self):
        yield self.token

    def terminals(self):
        yield self


class Edge:
    def __init__(self, pos, left, right=None, index=0, children=None):
        self.pos = pos
        self.left = left
        self.right = right or left
        self.index = index
        self.children = children or []

    def __repr__(self):
        str = []
        for i in range(len(self.pos.neighbors)):
            if i == self.index:
                str.append('^')
            str.append(self.pos.neighbors[i] + ' ')
        return '<%s => %s at %s:%s>' % (self.pos.label, "".join(str)[:-1], self.left, self.right)

    def active(self):
        return self.index < len(self.pos.neighbors)


class Chart:

    def __init__(self, rules=[]):
        self._rules = rules
        self.num = 0

    def tags(self, s):
        yield s.token_type
        yield s
        yield Token
        if isinstance(s, Token):
            if s.string=="made":
                yield "VERB"
            if s.string=="I":
                yield "NOUN"
        if s=="I":
            yield "SUBJECTIVE_NOUN"
        elif s in ["cooked", "made"]:
            yield "VERB"
        else:
            yield "NOUN"
            #yield s.upper()
            #if s in ["olive", "red"]:
            #    yield "ADJECTIVE"

    def tokenize(self, s):
        for token in s.split(" "):
            yield token

    def parse(self, string_or_tokens):
        if isinstance(string_or_tokens, str):
            tokens = list(self.tokenize(string_or_tokens))
        else:
            tokens = list(string_or_tokens)
        n = len(tokens)
        if n>0:
            self.n = n
            self.edges = []
            self.constituents = []
            for i in range(n):
                self.edges.append(list())
                self.constituents.append(list())

            for i in range(n):
                token = tokens[i]
                self._add_token(token, i)
            for c in self.constituents[0]: 
                assert c.left==0
                if c.right==self.n:
                    yield c
    
    def _add_token(self, token, position):
        for tag in self.tags(token):
            self._add_constituent(PreTerminal(tag, token, position))

    def _add_constituent(self, constituent):
        self.constituents[constituent.left].append(constituent)
        for edge in self.edges[constituent.left]:
            self._advance_over(edge, constituent)
        for pos in self._rules:
            if pos.matches(constituent.label):
                edge = Edge(pos, constituent.left)
                self._advance_over(edge, constituent)

    def _advance_over(self, edge, constituent):
        pos = edge.pos
        if edge.right == constituent.left and pos.neighbors[edge.index] == constituent.label:
            self._add_edge(Edge(pos, edge.left, constituent.right, edge.index + 1, edge.children + [constituent]))

    def _add_edge(self, edge):
        self.num += 1
        if edge.active():
            if edge.right < self.n:
                self.edges[edge.right].append(edge)
                for constituent in self.constituents[edge.right]:
                    self._advance_over(edge, constituent)
        else:
            self._add_constituent(Constituent(edge.pos.label, edge.children, edge.left, edge.right))
