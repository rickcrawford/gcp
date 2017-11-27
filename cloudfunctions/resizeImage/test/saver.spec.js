'use strict'

const saver = require('../lib/saver')
const expect = require('chai').expect

describe('saver module', () => {
  
	describe('"save"', () => {
		it('should export a function', () => {
			expect(saver.save).to.be.a('function')
		})
	});


	describe('"save(url)"', () => {
		it('should be success', (done) => {
			saver.save('https://www.onemarketnetwork.com/images/onemarket/content/team/michael-blandina.jpg')
			.then(() => done())
			.catch(err => {
				console.log("error", err);
			})
		});

		it('should be 304', (done) => {
			saver.save({url:'https://www.onemarketnetwork.com/images/onemarket/content/team/michael-blandina.jpg', etag: 'W/"wzwvixwhrspzezu6lc4dja=="'})
			.then(() => done())
			.catch(err => {
				console.log("error", err);
			})

		});


		it('should be 404', (done) => {
			saver.save('https://www.onemarketnetwork.com/images/onemarket/content/team/111.jpg')
			.catch(err => {
				done();
			});
		});
	});

});