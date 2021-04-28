"use strict";

function runTimer(timeZoneOffset) {
    setTimeout(() => {
        updateTime(timeZoneOffset)
    }, 1000)
    updateTime(timeZoneOffset)
}

function updateTime(timeZoneOffset) {
    let date = new Date()
    let gwtDate = new Date(date.getTime() + date.getTimezoneOffset() * 60000)
    let serverTime = new Date(gwtDate.getTime() + timeZoneOffset * 3600000)

    document.querySelector("#time").innerHTML =
        String(serverTime.getHours()).padStart(2, "0") + ":" +
        String(serverTime.getMinutes()).padStart(2, "0")
}