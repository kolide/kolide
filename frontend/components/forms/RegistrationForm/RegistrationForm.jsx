import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';

import AdminDetails from 'components/forms/RegistrationForm/AdminDetails';
import ConfirmationPage from 'components/forms/RegistrationForm/ConfirmationPage';
import KolideDetails from 'components/forms/RegistrationForm/KolideDetails';
import OrgDetails from 'components/forms/RegistrationForm/OrgDetails';

const PAGE_HEADER_TEXT = {
  1: 'SET USERNAME & PASSWORD',
  2: 'SET ORGANIZATION DETAILS',
  3: 'SET KOLIDE WEB ADDRESS',
  4: 'SUCCESS',
};

const baseClass = 'user-registration';

class RegistrationForm extends Component {
  static propTypes = {
    onNextPage: PropTypes.func,
    onSubmit: PropTypes.func,
    page: PropTypes.number,
  };

  constructor (props) {
    super(props);

    this.state = { errors: {}, formData: {} };
  }

  onPageFormSubmit = (pageFormData) => {
    const { formData } = this.state;
    const { onNextPage } = this.props;

    this.setState({
      formData: {
        ...formData,
        ...pageFormData,
      },
    });

    return onNextPage();
  }

  onSubmitConfirmation = () => {
    const { formData } = this.state;
    const { onSubmit: handleSubmit } = this.props;

    return handleSubmit(formData);
  }

  renderDescription = () => {
    const { page } = this.props;

    if (page === 1) {
      return (
        <div className={`${baseClass}__description`}>
          <p>Additional admins can be designated within the Kolide App</p>
          <p>Passwords must include 7 characters, at least 1 number (eg. 0-9) and at least 1 symbol (eg. ^&*#)</p>
        </div>
      );
    }

    if (page === 2) {
      return (
        <div className={`${baseClass}__description`}>
          <p>Set your Organization&apos;s name (eg. Yahoo! Inc)</p>
          <p>Specify the website URL of your organization (eg. Yahoo.com)</p>
        </div>
      );
    }

    if (page === 3) {
      return (
        <div className={`${baseClass}__description`}>
          <p>Define the base URL which osqueryd clients use to connect and register with Kolide.</p>
          <p>
            <small>Note: Please ensure the URL you choose is accessible to all endpoints that need to communicate with Kolide. Otherwise, they will not be able to correctly register.</small>
          </p>
        </div>
      );
    }

    return false;
  }

  renderHeader = () => {
    const { page } = this.props;
    const headerText = PAGE_HEADER_TEXT[page];

    if (headerText) {
      return <h2 className={`${baseClass}__title`}>{headerText}</h2>;
    }

    return false;
  }

  renderContent = () => {
    const { page } = this.props;
    const { formData } = this.state;
    const {
      onSubmitConfirmation,
      renderDescription,
      renderHeader,
    } = this;

    if (page === 4) {
      return (
        <div>
          {renderHeader()}
          <ConfirmationPage formData={formData} handleSubmit={onSubmitConfirmation} className={`${baseClass}__confirmation`} />
        </div>
      );
    }

    return (
      <div>
        {renderHeader()}
        {renderDescription()}
      </div>
    );
  }

  render () {
    const { onSubmit, page } = this.props;
    const { formData } = this.state;
    const { onPageFormSubmit, renderContent } = this;

    const containerClass = classnames(`${baseClass}__container`, {
      [`${baseClass}__container--complete`]: page > 3,
    });

    const adminDetailsClass = classnames(
      `${baseClass}__field-wrapper`,
      `${baseClass}__field-wrapper--admin`
    );

    const orgDetailsClass = classnames(
      `${baseClass}__field-wrapper`,
      `${baseClass}__field-wrapper--org`
    );

    const kolideDetailsClass = classnames(
      `${baseClass}__field-wrapper`,
      `${baseClass}__field-wrapper--kolide`
    );

    const formSectionClasses = classnames(
      `${baseClass}__form`,
      {
        [`${baseClass}__form--step1-active`]: page === 1,
        [`${baseClass}__form--step1-complete`]: page > 1,
        [`${baseClass}__form--step2-active`]: page === 2,
        [`${baseClass}__form--step2-complete`]: page > 2,
        [`${baseClass}__form--step3-active`]: page === 3,
        [`${baseClass}__form--step3-complete`]: page > 3,
      }
    );

    return (
      <div className={baseClass}>
        <div className={containerClass}>
          {renderContent()}

          <form onSubmit={onSubmit} className={formSectionClasses}>
            <AdminDetails formData={formData} handleSubmit={onPageFormSubmit} className={adminDetailsClass} />

            <OrgDetails formData={formData} handleSubmit={onPageFormSubmit} className={orgDetailsClass} />

            <KolideDetails formData={formData} handleSubmit={onPageFormSubmit} className={kolideDetailsClass} />
          </form>
        </div>
      </div>
    );
  }
}

export default RegistrationForm;
