async function GithubStats(username) {
    const GITHUB_API_URL = 'https://api.github.com';
    const LANGUAGE_TO_COLOR = {
        '1C Enterprise': '#814CCC',
        ABAP: '#E8274B',
        ActionScript: '#882B0F',
        Ada: '#02f88c',
        Agda: '#315665',
        'AGS Script': '#B9D9FF',
        Alloy: '#64C800',
        AMPL: '#E6EFBB',
        AngelScript: '#C7D7DC',
        ANTLR: '#9DC3FF',
        'API Blueprint': '#2ACCA8',
        APL: '#5A8164',
        AppleScript: '#101F1F',
        Arc: '#aa2afe',
        ASP: '#6a40fd',
        AspectJ: '#a957b0',
        Assembly: '#6E4C13',
        ATS: '#1ac620',
        AutoHotkey: '#6594b9',
        AutoIt: '#1C3552',
        Ballerina: '#FF5000',
        Batchfile: '#C1F12E',
        BlitzMax: '#cd6400',
        Boo: '#d4bec1',
        Brainfuck: '#2F2530',
        C: '#555555',
        'C#': '#178600',
        'C++': '#f34b7d',
        Ceylon: '#dfa535',
        Chapel: '#8dc63f',
        Cirru: '#ccccff',
        Clarion: '#db901e',
        Clean: '#3F85AF',
        Click: '#E4E6F3',
        Clojure: '#db5855',
        CoffeeScript: '#244776',
        ColdFusion: '#ed2cd6',
        'Common Lisp': '#3fb68b',
        'Common Workflow Language': '#B5314C',
        'Component Pascal': '#B0CE4E',
        Crystal: '#776791',
        CSS: '#563d7c',
        Cuda: '#3A4E3A',
        D: '#ba595e',
        Dart: '#00B4AB',
        DataWeave: '#003a52',
        DM: '#447265',
        Dockerfile: '#0db7ed',
        Dogescript: '#cca760',
        Dylan: '#6c616e',
        E: '#ccce35',
        eC: '#913960',
        ECL: '#8a1267',
        Eiffel: '#946d57',
        Elixir: '#6e4a7e',
        Elm: '#60B5CC',
        'Emacs Lisp': '#c065db',
        EmberScript: '#FFF4F3',
        EQ: '#a78649',
        Erlang: '#B83998',
        'F#': '#b845fc',
        Factor: '#636746',
        Fancy: '#7b9db4',
        Fantom: '#14253c',
        FLUX: '#88ccff',
        Forth: '#341708',
        Fortran: '#4d41b1',
        FreeMarker: '#0050b2',
        Frege: '#00cafe',
        'Game Maker Language': '#8fb200',
        Genie: '#fb855d',
        Gherkin: '#5B2063',
        Glyph: '#e4cc98',
        Gnuplot: '#f0a9f0',
        Go: '#375eab',
        Golo: '#88562A',
        Gosu: '#82937f',
        'Grammatical Framework': '#79aa7a',
        Groovy: '#e69f56',
        Hack: '#878787',
        Harbour: '#0e60e3',
        Haskell: '#5e5086',
        Haxe: '#df7900',
        HiveQL: '#dce200',
        HTML: '#e34c26',
        Hy: '#7790B2',
        IDL: '#a3522f',
        Idris: '#b30000',
        Io: '#a9188d',
        Ioke: '#078193',
        Isabelle: '#FEFE00',
        J: '#9EEDFF',
        Java: '#b07219',
        JavaScript: '#f1e05a',
        Jolie: '#843179',
        JSONiq: '#40d47e',
        Julia: '#a270ba',
        'Jupyter Notebook': '#DA5B0B',
        Kotlin: '#F18E33',
        KRL: '#28431f',
        Lasso: '#999999',
        Lex: '#DBCA00',
        LFE: '#4C3023',
        LiveScript: '#499886',
        LLVM: '#185619',
        LOLCODE: '#cc9900',
        LookML: '#652B81',
        LSL: '#3d9970',
        Lua: '#000080',
        Makefile: '#427819',
        Mask: '#f97732',
        Matlab: '#e16737',
        Max: '#c4a79c',
        MAXScript: '#00a6a6',
        Mercury: '#ff2b2b',
        Meson: '#007800',
        Metal: '#8f14e9',
        Mirah: '#c7a938',
        MQL4: '#62A8D6',
        MQL5: '#4A76B8',
        MTML: '#b7e1f4',
        NCL: '#28431f',
        Nearley: '#990000',
        Nemerle: '#3d3c6e',
        nesC: '#94B0C7',
        NetLinx: '#0aa0ff',
        'NetLinx+ERB': '#747faa',
        NetLogo: '#ff6375',
        NewLisp: '#87AED7',
        Nextflow: '#3ac486',
        Nim: '#37775b',
        Nit: '#009917',
        Nix: '#7e7eff',
        Nu: '#c9df40',
        'Objective-C': '#438eff',
        'Objective-C++': '#6866fb',
        'Objective-J': '#ff0c5a',
        OCaml: '#3be133',
        Omgrofl: '#cabbff',
        ooc: '#b0b77e',
        Opal: '#f7ede0',
        Oxygene: '#cdd0e3',
        Oz: '#fab738',
        P4: '#7055b5',
        Pan: '#cc0000',
        Papyrus: '#6600cc',
        Parrot: '#f3ca0a',
        Pascal: '#E3F171',
        PAWN: '#dbb284',
        Pep8: '#C76F5B',
        Perl: '#0298c3',
        'Perl 6': '#0000fb',
        PHP: '#4F5D95',
        PigLatin: '#fcd7de',
        Pike: '#005390',
        PLSQL: '#dad8d8',
        PogoScript: '#d80074',
        PostScript: '#da291c',
        PowerBuilder: '#8f0f8d',
        PowerShell: '#012456',
        Processing: '#0096D8',
        Prolog: '#74283c',
        'Propeller Spin': '#7fa2a7',
        Puppet: '#302B6D',
        PureBasic: '#5a6986',
        PureScript: '#1D222D',
        Python: '#3572A5',
        q: '#0040cd',
        QML: '#44a51c',
        R: '#198CE7',
        Racket: '#22228f',
        Ragel: '#9d5200',
        RAML: '#77d9fb',
        Rascal: '#fffaa0',
        Rebol: '#358a5b',
        Red: '#f50000',
        'Ren\'Py': '#ff7f7f',
        Ring: '#0e60e3',
        Roff: '#ecdebe',
        Rouge: '#cc0088',
        Ruby: '#701516',
        RUNOFF: '#665a4e',
        Rust: '#dea584',
        SaltStack: '#646464',
        SAS: '#B34936',
        Scala: '#c22d40',
        Scheme: '#1e4aec',
        sed: '#64b970',
        Self: '#0579aa',
        Shell: '#89e051',
        Shen: '#120F14',
        Slash: '#007eff',
        Smalltalk: '#596706',
        Solidity: '#AA6746',
        SourcePawn: '#5c7611',
        SQF: '#3F3F3F',
        Squirrel: '#800000',
        'SRecode Template': '#348a34',
        Stan: '#b2011d',
        'Standard ML': '#dc566d',
        SuperCollider: '#46390b',
        Swift: '#ffac45',
        SystemVerilog: '#DAE1C2',
        Tcl: '#e4cc98',
        Terra: '#00004c',
        TeX: '#3D6117',
        'TI Program': '#A0AA87',
        Turing: '#cf142b',
        TypeScript: '#2b7489',
        UnrealScript: '#a54c4d',
        Vala: '#fbe5cd',
        VCL: '#0298c3',
        Verilog: '#b2b7f8',
        VHDL: '#adb2cb',
        'Vim script': '#199f4b',
        'Visual Basic': '#945db7',
        Volt: '#1F1F1F',
        Vue: '#2c3e50',
        wdl: '#42f1f4',
        WebAssembly: '#04133b',
        wisp: '#7582D1',
        X10: '#4B6BEF',
        xBase: '#403a40',
        XC: '#99DA07',
        XQuery: '#5232e7',
        XSLT: '#EB8CEB',
        Yacc: '#4B6C4B',
        Zephir: '#118f9e'
    };

    /* We cache http response to prevent reaching rate limit & to load the page faster */
    function getCache(key) {
        return localStorage.getItem(`github-stats-${key}`);
    }

    function saveCache(key, value) {
        localStorage.setItem(`github-stats-${key}`, value);
    }

    function getJSONCache(key) {
        let value = getCache(key);
        return value ? JSON.parse(value) : null;
    }

    function saveJSONCache(key, value) {
        saveCache(key, JSON.stringify(value));
    }

    function getJSON(url) {
        return fetch(url)
            .then(response => {
                if (response.status !== 200) {
                    return;
                }
                return response.json();
            });
    }

    function getText(url) {
        return fetch(url)
            .then(response => {
                if (response.status !== 200) {
                    return;
                }
                return response.text();
            });
    }

    function countBytesWrittenInLanguage() {
        let languageCounts = {};
        repoLanguages.forEach(repoLanguage =>
            Object.keys(repoLanguage)
                .forEach(language =>
                    languageCounts[language] = (languageCounts[language] || 0) + repoLanguage[language]));
        return languageCounts;
    }

    function isUsernameChanged(username) {
        return getCache('username') !== username;
    }

    let cachedHomepageHTML = getCache('homepageHTML');
    let homepageHTML = !isUsernameChanged(username) && cachedHomepageHTML ?
        cachedHomepageHTML :
        await getText(`https://urlreq.appspot.com/req?method=GET&url=https://github.com/${username}`);
    saveCache('homepageHTML', homepageHTML);

    let parser = new DOMParser();
    let doc = parser.parseFromString(homepageHTML, 'text/html');
    let calendarSvg = doc.querySelector('.js-calendar-graph-svg');
    let days = calendarSvg.querySelectorAll('g > g > rect');

    Array.from(days).map(day => {
        return ({
            numCommits: day.dataset.count,
            date: day.dataset.date
        })
    });

    let commitsByDay = Array.from(days).map(day => {
        return ({
            numCommits: day.dataset.count,
            date: day.dataset.date
        })
    });

    let cachedRepos = getJSONCache('repos');
    let repos = !isUsernameChanged(username) && cachedRepos ?
        cachedRepos :
        await getJSON(`${GITHUB_API_URL}/users/${username}/repos?per_page=100`);

    let percentages;
    let repoLanguages;
    let languageCounts;
    let sortedLanguages;

    if(repos) {
        saveJSONCache('repos', repos);

        let cachedRepoLanguages = getJSONCache('repoLanguages');
        repoLanguages = !isUsernameChanged(username) && cachedRepoLanguages ?
            cachedRepoLanguages :
            await Promise.all(repos
                .filter(repo => !repo['private'])
                .map(repo =>
                    getJSON(`${GITHUB_API_URL}/repos/${repo['full_name']}/languages`)));

        if(repoLanguages) {
            repoLanguages = repoLanguages.filter(repoLanguage => repoLanguage);
            saveJSONCache('repoLanguages', repoLanguages);

            languageCounts = countBytesWrittenInLanguage(repoLanguages);

            let maxBytes = 0;
            for (let language in languageCounts) {
                if (languageCounts.hasOwnProperty(language) && languageCounts[language] > maxBytes) {
                    maxBytes = languageCounts[language];
                }
            }

            sortedLanguages = Object.keys(languageCounts).sort((languageA, languageB) => languageCounts[languageB] - languageCounts[languageA]);
            percentages = sortedLanguages.map(language => (languageCounts[language] / maxBytes));
        }
    }

    saveCache('username', username);

    return {
        commitsContribSVG: (config = {}) => {
            const rows = config.rows || 7;
            const space = config.space || 2;
            const rectWidth = config.rectWidth || 12;

            const levelColors = config.levelColors || [
                {
                    minCommits: 0,
                    color: '#ebedf0'
                },
                {
                    minCommits: 1,
                    color: '#c6e48b'
                },
                {
                    minCommits: 9,
                    color: '#7bc96f'
                },
                {
                    minCommits: 17,
                    color: '#239a3b'
                },
                {
                    minCommits: 26,
                    color: '#196127'
                }
            ];

            const cols = Math.ceil(commitsByDay.length / rows);

            let svgEl = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
            svgEl.setAttribute('width', `${cols * (rectWidth + space) - space}`);
            svgEl.setAttribute('height', `${rows * (rectWidth + space) - space}`);

            let groupEl = document.createElementNS('http://www.w3.org/2000/svg', 'g');

            for (let i = 0; i < commitsByDay.length; i++) {
                let col = Math.floor(i / rows);
                let row = i % rows;
                let rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                rect.setAttribute('x', `${col * (rectWidth + space)}`);
                rect.setAttribute('y', `${row * (rectWidth + space)}`);
                rect.setAttribute('width', `${rectWidth}`);
                rect.setAttribute('height', `${rectWidth}`);

                let color = null;
                let minCommits = -1;
                let mumCommits = commitsByDay[i].numCommits;

                levelColors.forEach(levelColor => {
                    if (mumCommits >= levelColor.minCommits && levelColor.minCommits > minCommits) {
                        minCommits = levelColor.minCommits;
                        color = levelColor.color;
                    }
                });

                rect.setAttribute('fill', color);
                groupEl.appendChild(rect);
            }

            svgEl.appendChild(groupEl);
            return svgEl;
        },
        languagesContribSVG: (config = {}) => {
            const barHeight = config.barHeight || 20;
            const barWidth =  config.barWidth || 500;

            const lineSpacing = config.lineSpacing || 4;
            const languageNameWidth = config.languageNameWidth || 100;
            const fontSize = config.fontSize || 14;

            let svgEl = document.createElementNS('http://www.w3.org/2000/svg', 'svg');

            if(sortedLanguages) {
                const svgHeight = (barHeight + lineSpacing) * sortedLanguages.length - lineSpacing;
                svgEl.setAttribute('width', `${barWidth + 100}`);
                svgEl.setAttribute('height', `${svgHeight}`);
                svgEl.setAttribute('viewport', `0 0 ${barWidth + 100} ${svgHeight}`);

                let languages = document.createElementNS('http://www.w3.org/2000/svg', 'g');
                languages.setAttribute('transform', `translate(${languageNameWidth},0)`);

                let bars = document.createElementNS('http://www.w3.org/2000/svg', 'g');
                bars.setAttribute('transform', `translate(${languageNameWidth},0)`);


                for (let i = 0; i < sortedLanguages.length; i++) {
                    let language = sortedLanguages[i];

                    let rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
                    let top = i * (barHeight + lineSpacing);

                    rect.setAttribute('y', `${top}`);
                    rect.setAttribute('width', `${percentages[i] * barWidth}`);
                    rect.setAttribute('height', `${barHeight}`);
                    rect.setAttribute('fill', LANGUAGE_TO_COLOR[language] || 'black');
                    bars.appendChild(rect);

                    let text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
                    text.setAttribute('font-size', fontSize);
                    text.style.fontWeight = 'bolder';
                    text.setAttribute('text-anchor', 'end');
                    text.setAttribute('dx', '-10');
                    text.setAttribute('y', `${barHeight / 2 + 4 + top}`);
                    text.textContent = `${language}`;
                    languages.appendChild(text);
                }
                svgEl.appendChild(languages);
                svgEl.appendChild(bars);
            }

            return svgEl;
        }
    };
}
