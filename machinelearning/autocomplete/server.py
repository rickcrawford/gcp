from flask import Flask, request, jsonify
import gensim, bisect, gzip, regex

def load_words(model):
	orig_words = [gensim.utils.to_unicode(word) for word in model.index2word]
	indices = [i for i, _ in sorted(enumerate(orig_words), key=lambda item: item[1].lower())]
	return [orig_words[i].lower() for i in indices]  # lowercased, sorted as lowercased


def related(term, model, all_words, confidence = 0.925, count = 10):
	results = []
	terms = []
	for t in term.split('_'):
		if t in all_words:
			terms.append(t)

	if len(terms) > 0:
		for t in model.most_similar(positive=terms, topn=count):
			if t[1] > confidence:
				results.append({"term": t[0], "confidence": t[1]})
	return results


app = Flask(__name__)

model3 = None
with gzip.open('model_trigram.bin.gz', 'rb') as f:
	model3 = gensim.models.KeyedVectors.load_word2vec_format(f, binary=True)
trigrams = load_words(model3)

# model2 = None
# with gzip.open('model_bigram.bin.gz', 'rb') as f:
# 	model2 = gensim.models.KeyedVectors.load_word2vec_format(f, binary=True)
# bigrams = load_words(model2)

model1 = None
with gzip.open('model_unigram.bin.gz', 'rb') as f:
	model1 = gensim.models.KeyedVectors.load_word2vec_format(f, binary=True)
unigrams = load_words(model1)



@app.route("/suggest", methods=['GET', 'POST'])
def suggest():
	

	if request.method == 'POST':
		prefix = request.form.get('prefix')
		confidence = request.form.get('confidence')
		count = request.form.get('count')
		rcount = request.form.get('rcount')
	else:
		prefix = request.args.get('prefix')
		confidence = request.args.get('confidence')
		count = request.args.get('count')
		rcount = request.args.get('rcount')

	prefix = gensim.utils.to_unicode(regex.sub("\\s+", "_", prefix)).strip().lower()

	if count is not None:
		count = int(count)

	if count is None or count < 0 or count > 20:
		count = 5

	if rcount is not None:
		rcount = int(rcount)

	if rcount is None or rcount < 0 or rcount > 20:
		rcount = 5

	if confidence is not None:
		confidence = float(confidence)

	if confidence is None or confidence >= 1.0 or confidence == 0:
		confidence = .95

	all_words = trigrams
	model = model1
	model_words = unigrams
	
	pos = bisect.bisect_left(all_words, prefix)
	words = all_words[pos: pos + count]

	results = []

	for word in words:
		if word.startswith(prefix):
			results.append({"term": ' '.join(word.split('_')), "related": related(word, model, model_words, confidence, rcount)})
	return jsonify({'data': results, 'meta':{'prefix':prefix, 'confidence': confidence, 'count':count, 'relatedCount': rcount}})


if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
