# import modules & set up logging
import gensim, logging, regex, bisect, gzip
logging.basicConfig(format='%(asctime)s : %(levelname)s : %(message)s', level=logging.INFO)

filename = 'products.txt.gz'

with gzip.open(filename, 'r') as f:
    content = f.readlines()

words = []
fixed = []
for line in content:
	for v in line.decode('utf-8').strip().lower().split(' - '):
		fixed.append(v)
		for word in regex.sub("[^\\p{L}\\p{N}]+", " ", v).split():
			words.append(word)

words = set(words) # so that all duplicate words are removed
fixed = set(fixed) # so that all duplicate words are removed

word2int = {}
int2word = {}
vocab_size = len(words) # gives the total number of unique words

print("Vocab Size: %d" % vocab_size)

for i,word in enumerate(words):
    word2int[word] = i
    int2word[i] = word

# raw sentences is a list of sentences.
sentences = []
for sentence in fixed:
    sentences.append(regex.split("[^\\p{L}\\p{N}]+", sentence))

print("Sentences: %d" % len(sentences))

# create phrases
bigram = gensim.models.Phrases(sentences)
trigram = gensim.models.Phrases(bigram[sentences])

print('Training word2vec models...')
# train word2vec on the two sentences
model_unigram = gensim.models.Word2Vec(sentences, size=200, window=4, min_count=5, workers=4)
model_unigram.init_sims(replace=True)

model_bigram = gensim.models.Word2Vec(bigram[sentences], size=100, min_count=2, workers=4)
model_bigram.init_sims(replace=True)

model_trigram = gensim.models.Word2Vec(trigram[bigram[sentences]], size=100, min_count=2, workers=4)
model_trigram.init_sims(replace=True)

# print('Saving artifacts to disk...')
import gzip

with gzip.open('model_unigram.bin.gz', 'wb') as f:
	model_unigram.wv.save_word2vec_format(f, binary=True)

with gzip.open('model_bigram.bin.gz', 'wb') as f:
	model_bigram.wv.save_word2vec_format(f, binary=True)

with gzip.open('model_trigram.bin.gz', 'wb') as f:
	model_trigram.wv.save_word2vec_format(f, binary=True)


# print(model.most_similar(positive=['microsoft','surface'], topn=10))


# with gzip.open('words.bin.gz', 'rb') as f:
# 	model = gensim.models.KeyedVectors.load_word2vec_format(f, binary=True)

# prefix = gensim.utils.to_unicode('micro').strip().lower()

# orig_words = [gensim.utils.to_unicode(word) for word in model.index2word]
# indices = [i for i, _ in sorted(enumerate(orig_words), key=lambda item: item[1].lower())]
# all_words = [orig_words[i].lower() for i in indices]  # lowercased, sorted as lowercased
# orig_words = [orig_words[i] for i in indices]  # original letter casing, but sorted as if lowercased

# count = 10
# pos = bisect.bisect_left(all_words, prefix)
# result = orig_words[pos: pos + count]
# print(result)
