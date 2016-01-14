import {GameTime} from './GameTime.js';

export class MessageLog
{
    constructor(maxMessageLines) {

        this.maxMessageLines = maxMessageLines;

        this.logMessages = [
            {time: 0, text: "Welcome to roguer!"}
        ];

        this.isDirty = false;

        var savedMessages = window.sessionStorage.getItem('_messages');
        if (savedMessages) {
            //console.log("restoring saved msgs from " + savedMessages);
            this.logMessages = JSON.parse(savedMessages);
        }
    }

    add(o)
    {
        if (!this.logMessages) {
            console.log("error: log wnd not yet ready!");
            console.log(o);
            return;
        }

        //console.log("added msg");
        //console.log(o);

        this.logMessages.push(o);
        this.isDirty = true;
        return this;
    }

    save()
    {
        if (!this.isDirty) {
            return;
        }

        this.isDirty = false;

        //console.log("saved message log in session storage");
        window.sessionStorage.setItem('_messages', JSON.stringify(this.logMessages));
    }

    render()
    {
        if (!this.logMessages) {
            console.log("error: log messages not yet ready!");
            return;
        }

        // TODO log window with scroll

        // only save the last messages in this.logMessages, and ignore scroll for now
        var msg = this.logMessages.slice(-this.maxMessageLines);

        var txt = "";
        for (let m of msg) {
            var time = new GameTime(m.time);
            txt = txt + time.render() + ": " + m.text + "\n";
        }

        return txt.trim();
    }
}
