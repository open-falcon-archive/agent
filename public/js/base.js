$(document).ready(function() {
 
    // enable popovers
    $(".pop").popover();

   // activate tooltips on hover
   $("[data-toggle='tooltip']").tooltip({trigger: 'hover', placement:'right'});  

    dashboard.getAll();
}).on("click", ".js-smoothscroll", function(event) {
    event.preventDefault();
    var target = $(this.hash).parent();
    pulseElement(target, 8, 400);

    $("html,body").animate({
        scrollTop: target.offset().top - 130
    }, 1000);
}).on("click", ".js-refresh-info", function(event) {
    event.preventDefault();
    var target = event.target;
    var item = target.id.split("-").splice(-1)[0];

    // if the refresh icon is click (where in a <span>) target will not have an id, so grab its parent instead
    if(target.id == "") {
        var parent = $(target).parent()[0];
        item = parent.id.split("-").splice(-1)[0];
    }

    dashboard.fnMap[item]();
});

// Handle for cancelling active effect.
var pulsing = {
    element: null,
    timeoutIDs: [],
    resetfn: function() {
        pulsing.element = null;
        pulsing.timeoutIDs = [];
    }
};

/**
 * Applies a pulse effect to the specified element. If triggered while already
 * active the ongoing effect is cancelled immediately.
 *
 * @param {HTMLElement} element The element to apply the effect to.
 * @param {Number} times How many pulses.
 * @param {Number} interval Milliseconds between pulses.
 */
function pulseElement(element, times, interval) {
    if (pulsing.element) {
        pulsing.element.removeClass("pulse").
            parent().removeClass("pulse-border");
        pulsing.timeoutIDs.forEach(function(ID) {
            clearTimeout(ID);
        });
        pulsing.timeoutIDs = [];
    }
    pulsing.element = element;
    var parent = element.parent();
    var f = function() {
        element.toggleClass("pulse");
        parent.toggleClass("pulse-border");
    };

    pulsing.timeoutIDs.push(setTimeout(pulsing.resetfn,
                                       (times + 1) * interval));
    for (; times > 0; --times) {
        pulsing.timeoutIDs.push(setTimeout(f, times * interval));
    }
}

function isInArray(array, search)
{
    return (array.indexOf(search) >= 0) ? true : false; 
}
