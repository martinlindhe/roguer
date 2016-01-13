import {GameTime, Hour} from './../GameTime.js';

describe('TimeOfDay', function () {
    it('shows time of day', function () {

        let t = new GameTime(Hour * 4);
        expect(t.render()).toEqual("04:00");
    });
});
