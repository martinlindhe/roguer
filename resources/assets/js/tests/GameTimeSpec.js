import {GameTime, Hour} from './../GameTime.js';

describe('GameTime', function () {
    it('shows time of day', function () {

        let t = new GameTime(Hour * 4);
        expect(t.render()).toEqual("04:00");
    });

    it('is day time', function () {

        let t1 = new GameTime(Hour * 4);
        expect(t1.isDaytime()).toEqual(false);

        let t2 = new GameTime(Hour * 15);
        expect(t2.isDaytime()).toEqual(true);
    });
});
