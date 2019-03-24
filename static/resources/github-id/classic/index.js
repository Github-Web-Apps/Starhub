$(document).ready(function () {
    let parse_query = function (query_string) {

        if (is_empty(query_string)) {
            return {};
        }
        let query = {};
        query_string = query_string.substring(query_string.indexOf('?') + 1);
        console.log(query_string);
        let query_array = query_string.split('&');
        query_array.forEach(function (q) {
            let kv = q.split('=');
            console.log(kv);
            query[kv[0]] = kv[1];
        });
        return query;

    };

    let is_empty = function (s) {
        return s === null || s === undefined || s === '';
    };

    let datetime_to_datetime_str = function (datetime) {
        let year = datetime.getFullYear();
        let month = left_padding_zero((datetime.getMonth() + 1));
        let day = left_padding_zero(datetime.getDate());
        let hour = left_padding_zero(datetime.getHours());
        let min = left_padding_zero(datetime.getMinutes());
        return year + '-' + month + '-' + day + ' ' + hour + ':' + min;
    };

    let left_padding_zero = function (int) {
        return int < 10 ? '0' + int : int;
    };

    let invoke_github_api = function (url, callback) {
        $.ajax({
            headers: {
                Accept: 'application/vnd.github.v3.star+json; charset=utf-8'
            },
            url: url,
            dataType: 'json',
            type: 'get',
            async: true,
            success: function (data, status, xhr) {
                callback(data, xhr);
            },
            error: function (xhr, status, error) {
                console.error(xhr.responseText);
                callback({}, xhr);
                if (xhr.status === 403) {
                    let query = parse_query(url);
                    console.log('qq', query);
                    let access_token = query['access_token'];
                    if (is_empty(access_token)) {
                        alert('Please wait for a while, printer is too hot');
                    } else {
                        $.ajax({
                            headers: {
                                Accept: 'application/vnd.github.v3.star+json; charset=utf-8'
                            },
                            url: 'https://api.github.com/rate_limit?access_token=' + access_token,
                            dataType: 'json',
                            type: 'get',
                            async: true,
                            success: function (data, status, xhr) {
                                alert('Rate limit exceeded, come back at ' + datetime_to_datetime_str(new Date(data['resources']['core']['reset'] * 1000)));
                            },
                            error: function (xhr, status, error) {
                                console.error(xhr.responseText);
                            }
                        });
                    }
                }
            }
        });
    };

    let random_int = function (start, end) {
        return Math.floor(start + (end - start) * Math.random());
    };

    let select_token = function () {
        const access_tokens = [
            'c36f15bcb6df04165eaa2f0b2aee3015eea499b0',
            '1b609d0389cd14e2b68446fdd09c1f37a783ee1e',
            '0d327a7f4ecc8758a2d092f5e33553816b0af761',
            '04b384b9e63f9a7c12f6407b3de0939b607d9d42',
            
            '4c61c7e58785bb8b2789224381b2e1331083ffa7',
            'a95aaa5a06cf5aa871ccc2f670371cfcb556cbd8',
            '46dcf5dead1cb95385eab0ebab674e3dd3601dcc',
            '6226fb196321a8a1b685b735e7208b909e12ef48',
            '76183e52d0b0b1e9265088be5197df9f1a14454e'];

        return access_tokens[random_int(0, access_tokens.length)];
    };

    let format_email = function (user_id, email) {
        if (is_empty((email))) {
            return '';
        } else if (user_id.length * 3 + email.length > 55) {
            let end = 55 - user_id.length * 3;
            if (end < 1) {
                return '...';
            } else {
                return email.substring(0, end) + '...';
            }
        } else {
            return email;
        }
    };

    let simple_number = function (num, digits) {
        if (num < 1000) {
            return num;
        } else if (1000 <= num && num < 1000000) {
            return (num % 1000 === 0 ? num / 1000 : (num / 1000).toFixed(digits)) + 'k';
        } else {
            return (num % 1000000 === 0 ? num / 1000000 : (num / 1000000).toFixed(digits)) + 'm';
        }
    };

    let object_to_array = function (object) {
        let language_array = [];
        Object.keys(object).forEach(function (k, i) {
            language_array.push([k, object[k]]);
        });
        return language_array;
    };

    let display_rank = function (data) {
        let myChart = echarts.init(document.getElementById('rank'));

        let option = {
            grid: {
                top: 0,
                bottom: 0,
                left: -16,
                right: -16
            },
            xAxis: {
                type: 'category',
                show: false
            },
            yAxis: {
                type: 'value',
                show: false

            },
            series: [{
                data: data,
                type: 'line',
                showSymbol: false,
                areaStyle: {},
                color: 'rgba(40,167,69,0.2)'

            }]
        };

        myChart.setOption(option);
    };

    let display_language = function (q, name, data) {
        let myChart = echarts.init(document.getElementById('language'));

        let option = {
            tooltip: {},
            angleAxis: {
                min: 0,
                max: 1,
                axisTick: {
                    show: false
                },
                axisLabel: {
                    show: true,
                    fontSize: 10,
                    margin: 5,
                    formatter: function (value, index) {
                        if (index === 1) {
                            return value * 100 + '%';
                        } else {
                            return '';
                        }
                    }
                },
                splitLine: {
                    lineStyle: {
                        type: 'dashed'
                    }
                }
            },
            radiusAxis: {
                type: 'category',
                data: name,
                z: 10,
                axisLabel: {
                    fontSize: 10,
                    interval: 0
                }
            },
            polar: {
                radius: '55%'
            },
            series: [{
                type: 'bar',
                data: data,
                coordinateSystem: 'polar',
                name: 'Language Percent',
                color: 'rgba(40,167,69,1.0)'
            }]
        };

        myChart.setOption(option);
    };

    let display_repo = function (data) {
        let myChart = echarts.init(document.getElementById('repo'));

        let option = {
            title: {
                subtext: 'Popular Repositories',
            },
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'shadow'
                }
            },
            grid: {
                top: 35,
                bottom: 20,
                left: 0,
                right: 15,
                containLabel: true
            },
            xAxis: {
                name: 'Stars',
                type: 'value',
                axisLabel: {
                    formatter: function (value, index) {
                        // 格式化成月/日，只在第一个刻度显示年份
                        if (index % 2 === 1) {
                            return simple_number(value, 1);
                        } else {
                            return '';
                        }
                    }
                },
                splitLine: {
                    lineStyle: {
                        type: 'dashed'
                    }
                }
            },
            yAxis: {
                type: 'category',
                axisLabel: {
                    show: false
                }

            },
            series: [
                {
                    name: 'Stars',
                    type: 'bar',
                    label: {
                        normal: {
                            show: true,
                            position: 'insideLeft',
                            formatter: function (param) {
                                return param['name'];
                            }
                        }
                    },
                    data: data
                },
            ],
            color: ['#0d7884']
        };

        myChart.setOption(option);
    };

    let render_chart = function (user_id) {
        let user_url = 'https://api.github.com/users/' + user_id + '?access_token=' + select_token();
        invoke_github_api(user_url, function (user_data) {
            let email = format_email(user_id, user_data['email']);
            let email_html = is_empty(email) ? '' : '<span style="font-size: 12px; color: #999; letter-spacing: 0.01em">&nbsp;' + email + '</span>';
            $('#user_id').html(user_id + email_html);
            $('#avatar').attr('src', user_data['avatar_url']);
            $('#following').html('Following: ' + simple_number(user_data['following'], 1));
            $('#follower').html('Follower: ' + simple_number(user_data['followers'], 1));

            let progress_bar = $('#progress-bar');
            progress_bar.css('width', '10%');

            let repo_count = user_data['public_repos'];
            if (repo_count === 0) {
                progress_bar.css('width', '100%');
                progress_bar.css('background-color', '#fff');
                return;
            }
            let page_size = 100;
            let page_count = Math.ceil(repo_count / page_size);
            let repos = [];
            for (let page = 1; page <= page_count; page++) {
                let repo_url = 'https://api.github.com/users/' + user_id + '/repos?sort=created&direction=asc&per_page='
                    + page_size + '&page=' + page + '&access_token=' + select_token();

                invoke_github_api(repo_url, function (repo_data) {
                    repos = repos.concat(repo_data.map(function (e) {
                        return [e['name'], e['stargazers_count']];
                    }));

                    progress_bar.css('width', (10 + 40 / repo_count * repos.length) + '%');

                    if (repos.length === repo_count) {
                        display_rank(repos);
                        repos.sort(function (r1, r2) {
                            return r2[1] - r1[1];
                        });

                        display_repo(repos.slice(0, 3).map(function (e) {
                            return [e[1], e[0]]
                        }).reverse());

                        let repos_no_io = repos.filter(function (r) {
                            return r[0].indexOf(user_id + '.github.io') === -1;
                        });

                        if (repos_no_io.length === 0) {
                            progress_bar.css('width', '100%');
                            progress_bar.css('background-color', '#fff');
                            return;
                        }

                        let language = {};
                        let load_repo_count = 0;
                        for (let i = 0; i < repos_no_io.length; i++) {
                            let r = repos_no_io[i];

                            let language_url = 'https://api.github.com/repos/' + user_id + '/' + r[0] + '/languages'
                                + '?access_token=' + select_token();
                            invoke_github_api(language_url, function (language_data, xhr) {
                                if (xhr.status === 200) {
                                    let total_codes = 0;
                                    Object.keys(language_data).forEach(function (k, i) {
                                        total_codes += language_data[k];
                                    });
                                    Object.keys(language_data).forEach(function (k, i) {
                                        let percent = language_data[k] / (total_codes * repos_no_io.length);
                                        if (language.hasOwnProperty(k)) {
                                            language[k] += percent;
                                        } else {
                                            language[k] = percent;
                                        }
                                    });
                                }

                                load_repo_count++;
                                progress_bar.css('width', (50 + 50 / repos_no_io.length * load_repo_count) + '%');

                                if (load_repo_count === repos_no_io.length) {

                                    let language_array = object_to_array(language);
                                    language_array.sort(function (l1, l2) {
                                        return l2[1] - l1[1];
                                    });

                                    let lang_name = [];
                                    let lang_data = [];
                                    language_array.slice(0, 5).forEach(function (l) {
                                        lang_name.push(l[0]);
                                        lang_data.push(l[1]);
                                    });

                                    if (lang_name.length > 0) {
                                        display_language('', lang_name.reverse(), lang_data.reverse());
                                    }

                                    progress_bar.css('background-color', '#fff');
                                }

                            })
                        }
                    }
                });
            }
        });
    };

    let query = parse_query(window.location.search);
    let q = is_empty(query['q']) ? 'intika' : query['q'];

    render_chart(q);

    $('#capture').click(function () {
        domtoimage.toPng(document.getElementById('outer'))
            .then(function (dataUrl) {
                let link = document.createElement('a');
                link.download = 'github-id_' + q + '.png';
                link.href = dataUrl;
                link.click();
            })
            .catch(function (error) {
                console.error('Oops, something went wrong!', error);
            });
    });

});