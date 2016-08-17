/**
 * @flow
 */

import React from 'react';
import Dispatcher from 'frontend/Dispatcher';

import { UserGetters } from 'frontend/stores/User';

import Login from 'frontend/components/Login';

function requireAuthentication(higherOrderComponent: $FlowFixMe) {
  return React.createClass({
    mixins: [Dispatcher.ReactMixin],

    getDataBindings() {
      return {
        logged_in: UserGetters.logged_in,
      }
    },

    getInitialState() {
      return {
        logged_in: null,
      }
    },

    componentWillMount() {
      this.checkAuth();
    },

    componentWillReceiveProps(nextProps: $FlowFixMe) {
      this.checkAuth();
    },

    checkAuth() {
      if (!this.state.logged_in) {
        // TODO: set redirect state so user can come back here after logging in
        this.props.router.push(Login.getRoute());
      }
    },

    render() {
      var Component = higherOrderComponent;
      return (
        <div>
          {this.state.logged_in === true
            ? <Component {...this.props}/> 
            : null
          }
        </div>
      )
    },
  })
}

exports.requireAuthentication = requireAuthentication