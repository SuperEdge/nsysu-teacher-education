const STATUS = {
    '': '未審核',
    'pass': '通過',
    'failed': '未通過',
}

$(document).ready(function () {
    $.ajax({
        url: config.server + '/v1/course/sign-up',
        type: 'GET',
        error: function (xhr) {
            console.error(xhr);
        },
        beforeSend: function (xhr) {
            let token = $.cookie('token')
            if (token == undefined) {
                renewToken()
                token = $.cookie('token')
            }

            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        },
        success: function (response) {
            if (response.list.length == 0) {
                $('#course tbody').append('\
                        <tr>\
                            <td scope="row" colspan="8" style="text-align: center">尚無資料</td>\
                        </tr>\
                    ')
            } else {
                response.list.forEach(function (element, index) {
                    let startDate = element.Course.Start.substring(0, 10)
                    let startTime = element.Course.Start.substring(11, 19)
                    let endDate = element.Course.End.substring(0, 10)
                    let endTime = element.Course.End.substring(11, 19)
                    let time = ""

                    if (startDate == endDate) {
                        time = startDate + ' ' + startTime + ' ~ ' + endTime
                    } else {
                        time = startDate + ' ' + startTime + ' ~ ' + endDate + ' ' + endTime
                    }

                    $('#student-course tbody').append('\
                        <tr>\
                            <th scope="row">'+ index + '</th>\
                            <td>'+ element.Course.Topic + '</td>\
                            <td>'+ time + '</td>\
                            <td>'+ element.Course.Type + '</td>\
                            <td>'+ STATUS[element.Status] + '</td>\
                            <td>'+ element.Comment + '</td>\
                            <td>'+ element.Review + '</td>\
                            <td class="edit"><button class="btn btn-primary" id="' + element.ID + '">編輯</button></td>\
                        </tr>\
                    ')
                })
            }
        }
    });
})

$('#student-course').on('click', '.edit', function () {
    let id = $(this).children().get(0).id
    let review = $(this).prev().html()

    $('#updateReviewModal textarea').val(review)
    $('#updateReviewModal input').val(id)
    $('#updateReviewModal').modal('show')
})

$('#updateReviewModal form').on('submit', function (e) {
    e.preventDefault()

    let studentCourseID = $('#updateReviewModal input').val()
    let review = $('#updateReviewModal textarea').val()

    $.ajax({
        url: config.server + '/v1/course/review',
        type: 'PATCH',
        error: function (xhr) {
            console.error(xhr);
        },
        data: {
            'StudentCourseID': studentCourseID,
            'Review': review,
        },
        beforeSend: function (xhr) {
            let token = $.cookie('token')
            if (token == undefined) {
                renewToken()
                token = $.cookie('token')
            }

            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        },
        success: function (response) {
            $('#updateReviewModal').modal('hide')
            $('button#' + studentCourseID).parent().prev().html(review)

            $('#alert').addClass('alert alert-success')
            $('#alert').html('成功<button type="button" class="close"><span aria-hidden="true">&times;</span></button>')
            $('#alert').show('fast')
        }
    })
})

$('#alert').on('click', '.close', function () {
    $('#alert').hide('slow')
})