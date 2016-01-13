// time constants
export const Minute = 1;
export const Hour   = Minute * 60;
export const Day    = Hour * 24;
export const Month  = Day * 30;
export const Season = Month * 3;
export const Year   = Month * 12;

export class GameTime
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

    minute()
    {
        var parts = this.dateParts();
        return parts.Minute;
    }

    hour()
    {
        var parts = this.dateParts();
        return parts.Hour;
    }

    dateParts()
    {
        var rest = this.time;

        var year = Math.floor(rest / Year);
        if (year > 0) {
            rest -= year * Year;
        }

        var month = Math.floor(rest / Month);
        if (month > 0) {
            rest -= month * Month;
        }

        var day = Math.floor(rest / Day);
        if (day > 0) {
            rest -= day * Day;
        }

        // NOTE: day is stored 0-based, but displayed 1-based
        day++;

        var hour = Math.floor(rest / Hour);
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
