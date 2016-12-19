import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import AceEditor from 'react-ace';
import { isEqual } from 'lodash';

import queryInterface from 'interfaces/query';
import Dropdown from 'components/forms/fields/Dropdown';

const baseClass = 'search-pack-query';

class SearchPackQuery extends Component {
  static propTypes = {
    allQueries: PropTypes.arrayOf(queryInterface),
    onSelectQuery: PropTypes.func,
    selectedQuery: queryInterface,
  };

  constructor (props) {
    super(props);

    this.state = {
      queryDropdownOptions: [],
    };
  }

  componentWillMount () {
    const { allQueries } = this.props;

    const queryDropdownOptions = allQueries.map((query) => {
      return { label: query.name, value: String(query.id) };
    });

    this.setState({ queryDropdownOptions });
  }

  componentWillReceiveProps (nextProps) {
    const { allQueries } = nextProps;

    if (!isEqual(allQueries, this.props.allQueries)) {
      const queryDropdownOptions = allQueries.map((query) => {
        return { label: query.name, value: String(query.id) };
      });

      this.setState({ queryDropdownOptions });
    }
  }

  renderHeader = () => {
    const { selectedQuery } = this.props;
    if (selectedQuery) {
      return (
        <h1 className={`${baseClass}__title`}>
          {selectedQuery.name}
        </h1>
      );
    }

    return (<h1 className={`${baseClass}__title`}>Choose Query</h1>);
  }

  renderQuery = () => {
    const { selectedQuery } = this.props;
    if (selectedQuery) {
      return (
        <AceEditor
          editorProps={{ $blockScrolling: Infinity }}
          mode="kolide"
          minLines={1}
          maxLines={3}
          name="pack-query"
          readOnly
          setOptions={{ wrap: true }}
          showGutter={false}
          showPrintMargin={false}
          theme="kolide"
          value={selectedQuery.query}
          width="100%"
          fontSize={14}
        />
      );
    }

    return false;
  }

  renderDescription = () => {
    const { selectedQuery } = this.props;
    if (selectedQuery) {
      return (
        <div className={`${baseClass}__description`}>
          <h2>description</h2>
          <p>{selectedQuery.description}</p>
        </div>
      );
    }

    return false;
  }

  render () {
    const { renderHeader, renderQuery, renderDescription } = this;
    const { onSelectQuery, selectedQuery } = this.props;
    const { queryDropdownOptions } = this.state;

    return (
      <div className={baseClass}>
        {renderHeader()}
        <Dropdown
          options={queryDropdownOptions}
          onChange={onSelectQuery}
        />
        {renderQuery()}
        {renderDescription()}
      </div>
    );
  }
}

export default SearchPackQuery;
