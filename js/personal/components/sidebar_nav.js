import css from './sidebar_nav.css';
import React, { Component, PropTypes } from 'react';
import { Link, withRouter, matchPath } from 'react-router-dom';
import style from 'js/common/hoc/style';
import { withStore } from 'js/personal/store';
import { observer } from 'mobx-react';
import cx from 'classnames';
import Menu from 'js/personal/components/menu';

const { object, bool } = PropTypes;

@withRouter
@withStore
@style(css)
@observer
export default class SidebarNav extends Component {
    static propTypes = {
        store: object.isRequired,
        visible: bool,
        location: object.isRequired,
    };

    constructor(...args) {
        super(...args);

        this.state = {
            showProjectsMenu: false,
        };
    }

    static defaultProps = {
        visible: true,
    };

    onProjMenuClick = () => {
        const { showProjectsMenu } = this.state;
        this.setState({ showProjectsMenu: !showProjectsMenu });
    }

    renderItem(project, path, title) {
        const { location } = this.props;
        const projectPath = `/projects/${ project.id }`;
        const selected = !!matchPath(location.pathname, `${ projectPath }${ path }`);

        return (
            <div className={ cx('sidebar-nav-item', { selected }) }>
                <Link to={ `${ projectPath }${ path }` }>{ title }</Link>
            </div>
        );
    }

    renderProjects() {
        const { store } = this.props;
        const { showProjectsMenu } = this.state;
        const projects = store.auth.projects.value;
        if (!projects) {
            return null;
        }
        if (showProjectsMenu) {
            return (
                <Menu className="projects-select" visible={ showProjectsMenu } onClose={ this.onProjMenuClick }>
                    {projects.map((p, i) => (
                        <Link key={ p.id } to={ `/projects/${ p.id }` }>{ p.name }</Link>
                  ))}</Menu>
            );
        }
    }

    render() {
        const { store, visible } = this.props;
        const activeProject = store.auth.activeProject;
        if (!activeProject) {
            return null;
        }

        const cls = cx('sidebar-nav', { visible });

        return (
            <div className={ cls }>
                <div className="project">
                    <div className="icon" style={ { backgroundImage: 'url(' + activeProject.icon_url + ')' } } />
                    <p onClick={ this.onProjMenuClick }>{ activeProject.name }</p>
                    { this.renderProjects() }
                </div>
                { this.renderItem(activeProject, '/conversations', 'Flows') }
                { this.renderItem(activeProject, '/faq', 'FAQ') }
                { this.renderItem(activeProject, '/assets', 'Assets') }
            </div>
        );
    }
}
