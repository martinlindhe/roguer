class MessageLog
{
    constructor() {
        this.logMessages = [];

        var savedMessages = window.sessionStorage.getItem('_messages');
        if (savedMessages) {
            console.log("restoring saved msgs from " + savedMessages);
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

        this.logMessages.push(o);
        return this;
    }

    save()
    {
        // console.log("saved message log in session storage");

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
        this.logMessages = this.logMessages.slice(-this.maxMessages);

        var txt = "";
        for (let msg of this.logMessages) {
            txt = txt + msg.time + ": " + msg.text + "\n";
        }

        return txt.trim();
    }
}

export default MessageLog;
