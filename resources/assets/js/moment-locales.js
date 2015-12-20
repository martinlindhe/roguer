import moment from 'moment';

moment.locale('en_US', {
    relativeTime : {
        future: "in %s",
        past:   "%s ago",
        s:  "seconds",
        m:  "a minute",
        mm: "%d minutes",
        h:  "an hour",
        hh: "%d hours",
        d:  "a day",
        dd: "%d days",
        M:  "a month",
        MM: "%d months",
        y:  "a year",
        yy: "%d years"
    }
});

moment.locale('sv_SE', {
    relativeTime : {
        future: "om %s",
        past:   "%s sedan",
        s:  "sekunder",
        m:  "en minut",
        mm: "%d minuter",
        h:  "en timme",
        hh: "%d timmar",
        d:  "en dag",
        dd: "%d dagar",
        M:  "en månad",
        MM: "%d månader",
        y:  "ett år",
        yy: "%d år"
    }
});
