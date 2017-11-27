'use strict'

const index = require('../index')
const expect = require('chai').expect

describe('index', () => {
  
	describe('"resizeImage"', () => {
		it('should export a function', () => {
			expect(index.resizeImage).to.be.a('function')
		})
	});

//http://img.bbystatic.com/BestBuy_US/images/products/3465/346575_rc.jpg
	describe('"resizeImage(data)"', () => {
		it('should be success', (done) => {

			let messageStr = new Buffer(JSON.stringify([
				{
					url: 'https://www.onemarketnetwork.com/images/onemarket/content/team/michael-blandina.jpg',
				},
				{
					url: 'https://www.onemarketnetwork.com/images/onemarket/content/team/don-kingsborough.jpg',
				}
			])).toString('base64')

			let message = {
				data: {
					data: messageStr
				}
			}

			index.resizeImage(message).then(() => {
				done();
			})
		});
	});
});