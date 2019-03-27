/*
    Simple ReactJS-based GitHub star browser.
    Sebastian Hanula - http://www.hanula.com

    Date: 2015-01-09
*/

function log() {
    if(window.console)
        return console.log.apply(console, arguments)
}

var LocalStorageMixin = {
    componentWillMount: function() {
        var key = this.storageKey;

        if(key && localStorage[key]) {
            this.setState(JSON.parse(localStorage[key]))
        }
    },
    componentDidUpdate: function(prevProps, prevState) {
        var key = this.storageKey;
        if(key) {
            localStorage[key] = JSON.stringify(this.state);
        }
    }
}


var Loading = React.createClass({
    render: function() {
        return (
            <div>
                <h3>Loading...</h3>
                {this.props.message}
            </div>
        );
    }
});

var Repo = React.createClass({
    render: function () {
        var info = this.props.info,
            full_name = info.full_name.split('/').join(' / ');
        return (
            <div className="repo">
                <h5><a target="_blank" href={info.html_url}>{full_name}</a></h5>
                <p>{info.description}</p>
                <ul className="list-inline">
                    <li><b>&#9733;</b>   <strong>{info.stargazers_count}</strong></li>
                    <li><b>&#x2646;</b> <strong>{info.forks}</strong></li>
                </ul>
            </div>
        );
    }
});


// Filtered list of repositories.
var RepoList = React.createClass({

    filterText: function(text) {
        this.setState({filterText: text.trim()});
    },

    getInitialState: function() {
        return {filterText: ''};
    },

    shouldShowRepo: function(info) {
        if(this.state.filterText) {
            var text = this.state.filterText,
                exp = new RegExp(text, 'i');

            if(exp.test(info.full_name) ||
               exp.test(info.description)) {
                return true;
            }
            return false;
        }
        return true;
    },

    handleSaveClick: function() {
        this.props.onSaveFilter(this.state.filterText);
    },

    handleClear: function() {
        this.props.onClearFilter();
        this.setState({filterText: ''});
    },

    render: function () {
        // filter repos
        var repos = this.props.repos.filter(function(info) {
            if(this.shouldShowRepo(info)) return info;
        }.bind(this));

        var filter;

        if(this.state.filterText)
            filter = <div>
                <h4 className="pull-left">Filtering "{this.state.filterText}"</h4>:
                    <button className="btn btn-default" onClick={this.handleSaveClick}>save</button>
                    <button className="btn btn-default" onClick={this.handleClear}>clear</button>
                </div>;

        return (
            <div>

            <div className="page-header">
                <h3>Starred repositories <small>{repos.length}</small></h3>
                {filter}
            </div>

            <div className="row">
            {repos.map(function(info) {
                return <div key={info.id} className="repo-card col-md-3">
                            <Repo info={info}/>
                        </div>;
            }.bind(this))}
            </div>
            </div>
        );
    }
});

var UserSession = React.createClass({
    getInitialState: function() {
        return {changing: false}
    },
    handleChange: function(event) {
        var username = this.refs.input.getDOMNode().value;
        this.props.onChange(username);
        event.preventDefault();
    },

    render: function() {
        return <form onSubmit={this.handleChange} className="navbar-form">
                <input type="text"
                    ref="input"
                    id="inputLogin"
                    className="form-control"
                    defaultValue={this.props.username}
                    placeholder="GitHub username"/>
                <input type="submit" className="btn btn-default" value="Change" id="changeLogin" />
                <button className="btn" onClick={this.props.onLogout}>Logout</button>
            </form>
    }
});


var Filters = React.createClass({
    mixins: [LocalStorageMixin],
    storageKey: 'Filters',

    getInitialState: function() {
        return {filters: []}
    },
    add: function(text) {
        if(this.state.filters.indexOf(text) < 0) {
            this.state.filters.push(text);
            this.setState({filters: this.state.filters});
        }
    },
    remove: function(text) {
        this.setState({filters: this.state.filters.filter(function(t) {
            return t != text.trim();
        })});
    },
    handleClick: function(text) {
        this.props.onClick(text);
    },
    handleRemove: function(text, event) {
        event.stopPropagation();
        this.remove(text);
    },
    render: function () {
        return (
            <ul className="filters nav">
                {this.state.filters.map(function(text) {
                    return <li key={text} className="filter">
                        <a onClick={this.handleClick.bind(this, text)}>
                            {text}
                            <button onClick={this.handleRemove.bind(this, text)} className="remove">
                            &#x00D7;
                            </button>
                        </a>
                    </li>;
                }.bind(this))}
            </ul>
        );
    }
});

var Nav = React.createClass({
    mixins: [LocalStorageMixin],
    storageKey: 'Nav',

    getInitialState: function() {
        return {search: ''}
    },
    componentDidMount: function() {

        this.props.onSearch(this.state.search);
        this.refs.search.getDOMNode().focus();
    },
    handleSearch: function(event) {
        this.props.onSearch(event.target.value);
        this.setState({search: event.target.value});
    },
    clearSearch: function() {
        this.setState({search: ''});
        this.props.onSearch(this.state.search);
        this.refs.search.getDOMNode().value = '';
    },

    render: function () {
        return (
            <nav className="navbar navbar-default navbar-static-top">
                <div className="container">
                    <div className="navbar-header">
                      <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                        <span className="sr-only">Toggle navigation</span>
                        <span className="icon-bar"></span>
                        <span className="icon-bar"></span>
                        <span className="icon-bar"></span>
                      </button>

                      <a href="#"  className="navbar-brand" onClick={this.clearSearch}>
                        StarHub
                      </a>

                    </div>

                    <div id="navbar" className="navbar-collapse collapse">

                    <form className="navbar-form navbar-left">
                        <input ref="search"
                            onChange={this.handleSearch}
                            type="text"
                            className="form-control"
                            placeholder="Search..."/>
                    </form>

                      <ul className="nav navbar-nav navbar-right">
                        <li>
                            <UserSession username={this.props.username}
                                onChange={this.props.onUserChange}
                                onLogout={this.props.onLogout}/>
                        </li>
                      </ul>
                    </div>
                </div>
            </nav>
        );
    }
});

var UserStars = React.createClass({
    mixins: [LocalStorageMixin],
    storageKey: 'UserStars',

    loadStarredRepos: function(username) {
        this.setState({loading: true, loadingMessage: ''});

        GitHub.getUserStars(username).then(function(data) {
            this.setState({loading: false, starredRepos: data});
        }.bind(this)).fail(function(err) {
            this.setState({
                loadingMessage: ("Cannot retrieve data. " +
                                "Probably an API limit, try again in a minute. (note that this feature in its current state may not work with huge starred list)")
            });
        }.bind(this));
    },

    getInitialState: function() {
        return {
            loading: false,
            username: this.props.username,
            starredRepos: [],
        }
    },

    componentDidMount: function() {
        if(!this.state.starredRepos.length ||
           this.state.username != this.props.username)
            this.loadStarredRepos(this.props.username);
        this.state.loading = false;

    },

    handleSearch: function(text) {
        if(this.refs.repos)
            this.refs.repos.filterText(text);
    },

    handleSaveFilter: function(text) {
        this.refs.filters.add(text);
    },

    handleClearFilter: function() {
        this.refs.nav.clearSearch();
    },

    handleUserchange: function(username) {
        this.setState({username: username});
        this.loadStarredRepos(username);
    },

    applyFilter: function(text) {
        this.refs.repos.filterText(text);
    },

    render: function () {
        return (
                <div>
                <Nav
                    ref="nav"
                    onSearch={this.handleSearch}
                    username={this.state.username}
                    onLogout={this.props.onLogout}
                    onUserChange={this.handleUserchange}
                    />
            <div className="container">
            <div className="row">
                <div className="main-content col-md-10">
                {this.state.loading ? <Loading message={this.state.loadingMessage} /> :
                    <RepoList
                        ref="repos"
                        repos={this.state.starredRepos}
                        username={this.props.username}
                        onSaveFilter={this.handleSaveFilter}
                        onClearFilter={this.handleClearFilter}/>
                }
                </div>
            </div>
            <br></br>
            <br></br>
            <br></br>

            </div></div>
        )
    }
});

var UserSelect = React.createClass({
    handleLogin: function() {
        this.props.onLogin(this.refs.login.getDOMNode().value);
    },
    getInitialState: function() {
        return {}
    },
    handleInput: function() {
        this.setState({showHint: true});
    },
    render: function () {
        return (
            <div className="container">

                <div className="login-form">
                    <h1>
                    	<small>
                        Starred Repositories
                    	</small>    
                    </h1>
                    <form onSubmit={this.handleLogin}>
                        <input type="text"
                            ref="login"
                            id="inputLoginMain"
                            onChange={this.handleInput}
                            className="form-control"
                            placeholder="GitHub username"/>
                        <input type="submit" className="form-control" value="Search" id="submitLogin" />
                        <div className="hint">{this.state.showHint ? "hit <enter> to continue": ''}</div>
                    </form>
                </div>
                
            </div>
        )
    }
});

var Stars = React.createClass({
    mixins: [LocalStorageMixin],
    storageKey: 'Stars',

    getInitialState: function() {
        return {username: ''}
    },

    handleLogin: function(username) {
        this.setState({username: username});
    },
    handleLogout: function(){
        this.setState({username: ''});
        localStorage.clear();
    },
    render: function () {
        if(this.state.username) {
            return <UserStars
                ref="app"
                username={this.state.username}
                onLogout={this.handleLogout}/>
        } else {
            return <UserSelect onLogin={this.handleLogin}/>
        }
    }
});


React.render(<Stars/>, document.getElementById('app'));
