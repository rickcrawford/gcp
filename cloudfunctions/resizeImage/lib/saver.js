'use strict'

const maxLength = 10000000;
const defaultHeight = 50;
const defaultWidth = 50
const defaultBackground = 'white';
const entityType = 'Image';

const err304 = '304'

const request = require('request');
const sharp = require('sharp');
const sha256 = require('js-sha256').sha256;
const Buffer = require('safe-buffer').Buffer;
const datastore = require('@google-cloud/datastore')();

function getDatastoreKey(urlHash) {
  return datastore.key([entityType, urlHash]);
}


/**
 * saver creates an image asset.
 *
 * @param {options} url to request.
 */
exports.save = function(options) {
  return new Promise((resolve, reject) => {
    let url = '';
    let etag = '';
    let lastModified = '';
    let height = defaultHeight;
    let width = defaultWidth;
    let background = defaultBackground;

    if (typeof options == 'string' || options instanceof String) {
      url = options;
    } else if (typeof options == 'object' && options != null) {
      url = options.url;
      etag = options.etag ? options.etag : '';
      lastModified = options.lastModified ? options.lastModified : '';
      height = options.height ? options.height : defaultHeight;
      width = options.width ? options.width : defaultWidth;
      background = options.background ? options.background : defaultBackground;
    } else {
      return reject(new Error('no url'));
    }

    let reqOptions = {
      url: url,
      headers: {
        'If-None-Match': etag,
        'If-Modified-Since': lastModified, 
      },
    };

    const resizer = sharp()
      .resize(height, width)
      .background(background)
      .png();

    let response = {
      hash: sha256(url),
      metadata: {
        url: url,
        height: height,
        width: width,
        background: background,
        created: new Date(),
      }
    }


    // make a request, verify the response is an image.
    let req = request(reqOptions);
    req.on('response', (resp) => {
      response.statusCode = resp.statusCode;
      switch(resp.statusCode) {
        case 200:
          break;

        case 304:
          resp.destroy(err304);
          return;

        default:
          resp.destroy(`bad status: ${resp.statusCode}`);
          return;
      }

      let length = parseInt(resp.headers['content-length']);
      if (length == 0 || length > maxLength) {
        resp.destroy(`invalid length: ${length}`);
        return;
      }
      
      let contentType = resp.headers['content-type'];
      if (contentType.split("/")[0].toLowerCase().trim() != "image") {
        resp.destroy(`not an image: ${contentType}`);
        return;
      }

      response.metadata.etag = resp.headers['etag'];
      response.metadata.lastModified = resp.headers['last-modified'];
    });

    // let originalChunks = [];
    // req.on('data', chunk => {
    //   originalChunks.push(chunk);
    // })

    // there was an error... blow up.
    req.on('error', (err) => {
      if (err != err304) {
        return reject(new Error(err));
      }    
      return resolve(response);
    });

    // resize the image via a stream pipe
    let chunks = [];
    let resized = req.pipe(resizer)
    resized.on('data', chunk => {
      chunks.push(chunk);
    });

    resized.on('end', () => {
      // response.image.data = Buffer.concat(originalChunks);
      // response.image.hash = sha256(response.image.data);
      response.metadata.content = `image/png;base64,${Buffer.concat(chunks).toString('base64')}`;
      const entity = {
        key: getDatastoreKey(response.hash),
        data: response.metadata,
        excludeFromIndexes: ['content'],
      };
      return resolve(datastore.save(entity));
    });

    req.end();

  })
};
