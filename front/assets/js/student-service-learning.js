const TYPE = {
    'both': '同時認列教育實習服務暨志工服務',
    'internship': '實習服務',
    'volunteer': '志工服務',
}
const STATUS = {
    '': '未審核',
    'pass': '通過',
    'failed': '未通過',
}

let studentServiceLearningID = undefined

$(document).ready(() => {
    loading()

    Promise.all([
        getStudentInformation(),
        getStudentServiceLearning(),
    ]).then(() => {
        unloading()
    }).catch(() => {
        setTimeout(() => {
            removeCookie()
        }, 1500)
    })
})

const getStudentInformation = () => {
    $.ajax({
        url: `${config.server}/v1/user`,
        type: 'GET',
        beforeSend: (xhr) => {
            setHeader(xhr)
        },
        error: (xhr) => {
            errorHandle(xhr, "錯誤")
        },
        success: (response) => {
            student = Object.assign({}, response.data)
            $('div.greeting').html(`Hi, ${student.Name}同學`)
        }
    });
}

const getStudentServiceLearning = () => {
    $.ajax({
        url: `${config.server}/v1/service-learning/student`,
        type: 'GET',
        beforeSend: (xhr) => {
            $('#student-service-learning tbody').html('')
            setHeader(xhr)
        },
        error: (xhr) => {
            errorHandle(xhr, "錯誤")
        },
        success: (response) => {
            if (response.list.length == 0) {
                $('#student-service-learning tbody').append(`
                    <tr>
                        <td scope="row" colspan="8" style="text-align: center">尚無資料</td>
                    </tr>
                `)
            } else {
                response.list.forEach((element, index) => {
                    let color = 'class="waiting"'
                    if (element.Status === 'pass') {
                        color = 'class="success"'
                    } else if (element.Status === 'failed') {
                        color = 'class="danger"'
                    }

                    let date = `${dayjs(element.ServiceLearning.Start).format('YYYY-MM-DD')} ~ ${dayjs(element.ServiceLearning.End).format('YYYY-MM-DD')}`
                    let result = `
                        <tr>
                            <td data-title="審核情況" ${color}><span>●</span>${STATUS[element.Status]}</td>
                            <td data-title="類別">${TYPE[element.ServiceLearning.Type]}</td>
                            <td data-title="服務內容說明">${element.ServiceLearning.Content}</td>
                            <td data-title="日期">${date}</td>
                            <td data-title="時段">${element.ServiceLearning.Session}</td>
                            <td data-title="時數">${element.ServiceLearning.Hours}</td>
                    `

                    if (element.Status !== 'pass') {
                        result = `${result}<td><a class="btn_table" onclick="edit(${element.ID})">編輯</a></td></tr>`
                    } else {
                        result = `${result}<td><a class="btn_table disabled">編輯</a></td></tr>`
                    }

                    $('#student-service-learning tbody').append(result)
                })
            }
        }
    });
}

const edit = id => {
    studentServiceLearningID = id
    $('#Modal').modal('show')
}

const browse_file = (id, value) => {
    let splitted = value.split('\\')
    $(`#${id}`).val(splitted[splitted.length - 1])
}

$('#Modal form').on('submit', (e) => {
    e.preventDefault();

    let reference = $('#reference').prop('files')[0]
    let review = $('#review').prop('files')[0]

    if (reference === undefined && review === undefined) {
        swal({
            title: '',
            text: '請至少上傳一個檔案',
            icon: "error",
            timer: 1000,
            buttons: false,
        })
        return
    }

    if (reference.name.length > 36 || review.name.length > 36) {
        swal({
            title: '',
            text: '檔名過長',
            icon: "error",
            timer: 1000,
            buttons: false,
        })
        return
    }

    let formData = new FormData()
    formData.append('StudentServiceLearningID', studentServiceLearningID)

    if (reference !== undefined) {
        formData.append('Reference', reference)
    }
    if (review !== undefined) {
        formData.append('Review', review)
    }

    $.ajax({
        url: `${config.server}/v1/service-learning/student`,
        type: 'PATCH',
        data: formData,
        cache: false,
        contentType: false,
        processData: false,
        beforeSend: (xhr) => {
            setHeader(xhr)
        },
        error: (xhr) => {
            errorHandle(xhr, "上傳失敗")
        },
        success: (response) => {
            if (response.code == 0) {
                swal({
                    title: '',
                    text: '成功',
                    icon: "success",
                    timer: 1500,
                    buttons: false,
                }).then(() => {
                    $('#Modal').modal('hide')
                })
            } else {
                swal({
                    title: '',
                    text: '發生錯誤，請聯絡管理員',
                    icon: "error",
                    timer: 1500,
                    buttons: false,
                })
            }
        },
        complete: () => {
            $('#Modal form')[0].reset()
        }
    });
})

const newServiceLearning = () => {
    $('#submit').click()
}

$('#form').on('submit', (e) => {
    e.preventDefault();

    let type = $("#type").val()
    let content = $("#content").val()
    let startDate = $("#start-date input").val();
    let endDate = $("#end-date input").val();
    let startTime = $("#start-time input").val();
    let endTime = $("#end-time input").val();
    let hours = $("#hours").val()

    $.ajax({
        url: `${config.server}/v1/service-learning`,
        type: 'POST',
        data: {
            'Type': type,
            'Content': content,
            'Start': startDate,
            'End': endDate,
            'Session': `${startTime} ~ ${endTime}`,
            'Hours': hours
        },
        beforeSend: (xhr) => {
            loading()
            setHeader(xhr)
        },
        error: (xhr) => {
            swal({
                title: '',
                text: '發生錯誤',
                icon: "error",
                timer: 1500,
                buttons: false,
            })
        },
        success: (response) => {
            if (response.code == 0) {
                swal({
                    title: '',
                    text: '成功',
                    icon: "success",
                    timer: 1500,
                    buttons: false,
                }).then(() => {
                    $('#Modal').modal('hide')
                })
            } else {
                swal({
                    title: '',
                    text: '發生錯誤',
                    icon: "error",
                    timer: 1500,
                    buttons: false,
                })
            }
        },
        complete: () => {
            $('#form')[0].reset()
            getStudentServiceLearning()
            unloading()
        }
    });
})

$("#start-date input").on('change', () => {
    $("#end-date input").attr('min', $("#start-date input").val())
})

$("#end-date input").on('change', () => {
    $("#start-date input").attr('max', $("#end-date input").val())
})

$("#start-time input").on('change', () => {
    $("#end-time input").attr('min', $("#start-time input").val())
})

$("#end-time input").on('change', () => {
    $("#start-time input").attr('max', $("#end-time input").val())
})