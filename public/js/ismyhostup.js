var months = [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec'
]

$(function() {
    $('[tooltop]').each(function() {
        $(this).popup({
            content: $(this).attr('tooltip'),
            position: 'top center'
        });
    });
    $('.ui.dropdown').dropdown();
    $('.ui.checkbox').checkbox();
});

$(formatDates);

function formatDates() {
    $('.date-format').each(function() {
        var date = new Date(parseInt($(this).text()) * 1000);
        var day = date.getDay();
        var month = months[date.getMonth()];
        var year = date.getFullYear();
        var hour = date.getHours();
        var minute = date.getMinutes();
        $(this).text(day + ' ' + month + ' ' + year + ', ' + hour + ':' + (minute < 10 ? '0' : '' ) + minute);
        $(this).removeClass('date-format');
    });
}
