import GameTime from './../GameTime.js';

describe('TimeOfDay', function () {
    it('shows time of day', function () {
        var t = new GameTime(0);
        expect(t.render()).toEqual("00:00");
    });
});
