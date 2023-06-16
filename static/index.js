// credit: https://stackoverflow.com/a/60884408/18965816
window.addEventListener('load', () => {
    let now = new Date();
    const monthLen = new Date(now.getFullYear(), now.getDate(), 0);
    if (monthLen > now.getDate() + 1) {
        now.setDate(now.getDate() + 1)
    } else {
        if (now.getMonth() == 11) {
            now.setFullYear(now.getFullYear() + 1);
            now.setMonth(0)
        } else {
            now.setMonth(now.getMonth() + 1)
        }
        now.setDate(1)
    }
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());

    /* remove second/millisecond if needed - credit ref. https://stackoverflow.com/questions/24468518/html5-input-datetime-local-default-value-of-today-and-current-time#comment112871765_60884408 */
    now.setMilliseconds(null)
    now.setSeconds(null)

    document.getElementById('date').value = now.toISOString().slice(0, -1);
});
