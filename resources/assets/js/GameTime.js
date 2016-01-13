// time constants
const Minute = 1;
const Hour   = Minute * 60;
const Day    = Hour * 24;
const Month  = Day * 30;
const Season = Month * 3;
const Year   = Month * 12;

export default class GameTime
{
    constructor(i)
    {
        this.set(i);
    }

    set(i)
    {
        this.time = i;
    }

    /**
     * @return string "14:30"
     */
    render()
    {
        var parts = this.dateParts();

        return ("00" + parts.Hour).slice(-2) + ":" + ("00" + parts.Minute).slice(-2);
    }

    dateParts()
    {
        var rest = this.time;

        var year = rest / Year;
        if (year > 0) {
            rest -= year * Year;
        }

        var month = rest / Month;
        if (month > 0) {
            rest -= month * Month;
        }

        var day = rest / Day;
        if (day > 0) {
            rest -= day * Day;
        }

        // NOTE: day is stored 0-based, but displayed 1-based
        day++;

        var hour = rest / Hour;
        if (hour > 0) {
            rest -= hour * Hour;
        }

        var minute = rest;

        return {
            Minute: minute,
            Hour: hour,
            Day: day,
            Month: month,
            Year: year
        };
    }
}
