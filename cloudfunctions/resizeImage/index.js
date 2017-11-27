'use strict';

const saver = require('./lib/saver');

/**
 * Background Cloud Function to be triggered by Cloud Storage.
 *
 * @param {object} event The Cloud Functions event.
 * @param {function} callback The callback function.
 */
exports.resizeImage =  (event) => {
  const image = event.data || {};
  const data = JSON.parse(image.data ? Buffer.from(image.data, 'base64').toString() : '[]');
  const tasks = data.map(options => {
  	return saver.save(options)
  })
  return Promise.all(tasks);
};
