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

    render()
    {
        // XXXX convert int to HH:MM in js


        return "hehe " + this.time;
    }
}
