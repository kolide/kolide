import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Checkbox from 'components/forms/fields/Checkbox';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import Icon from 'components/Icon';
import InputField from 'components/forms/fields/InputField';
import Slider from 'components/forms/fields/Slider';
import validate from 'components/forms/admin/AppConfigForm/validate';

const authMethodOptions = [
  { label: 'Plain', value: 'plain' },
  { label: 'Login', value: 'login' },
  { label: 'GSS API', value: 'gssapi' },
  { label: 'Digest MD5', value: 'digest_md5' },
  { label: 'MD5', value: 'md5' },
  { label: 'Cram MD5', value: 'cram_md5' },
];
const authTypeOptions = [
  { label: 'Username and Password', value: 'username_and_password' },
  { label: 'None', value: 'none' },
];
const baseClass = 'app-config-form';
const formFields = [
  'auth_method', 'authentication_type', 'domain', 'enable_ssl_tls', 'enable_start_tls',
  'kolide_server_url', 'org_logo_url', 'org_name', 'password', 'port', 'sender_address',
  'server', 'user_name', 'verify_ssl_certs',
];
const Header = ({ showAdvancedOptions }) => {
  const CaratIcon = <Icon name={showAdvancedOptions ? 'downcarat' : 'upcarat'} />;

  return <span>Advanced Options {CaratIcon} <small>You normally don’t need to change these settings they are for special setups.</small></span>;
};

Header.propTypes = { showAdvancedOptions: PropTypes.bool.isRequired };

class AppConfigForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      auth_method: formFieldInterface.isRequired,
      authentication_type: formFieldInterface.isRequired,
      domain: formFieldInterface.isRequired,
      enable_ssl_tls: formFieldInterface.isRequired,
      enable_start_tls: formFieldInterface.isRequired,
      kolide_server_url: formFieldInterface.isRequired,
      org_logo_url: formFieldInterface.isRequired,
      org_name: formFieldInterface.isRequired,
      password: formFieldInterface.isRequired,
      port: formFieldInterface.isRequired,
      sender_address: formFieldInterface.isRequired,
      server: formFieldInterface.isRequired,
      user_name: formFieldInterface.isRequired,
      verify_ssl_certs: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    smtpConfigured: PropTypes.bool,
  };

  constructor (props) {
    super(props);

    this.state = { showAdvancedOptions: false };
  }

  onToggleAdvancedOptions = (evt) => {
    evt.preventDefault();

    const { showAdvancedOptions } = this.state;

    this.setState({ showAdvancedOptions: !showAdvancedOptions });

    return false;
  }

  renderAdvancedOptions = () => {
    const { fields } = this.props;
    const { showAdvancedOptions } = this.state;

    if (!showAdvancedOptions) {
      return false;
    }

    return (
      <div>
        <div className={`${baseClass}__inputs`}>
          <div className={`${baseClass}__smtp-section`}>
            <InputField {...fields.domain} label="Domain" />

            <div className="input-field__wrapper">
              <label
                className="input-field__label"
                htmlFor={fields.verify_ssl_certs.name}
              >
                Verify SSL Certs?
              </label>
              <div className="slide-wrapper">
                <span className="slider-option slider-option--off">OFF</span>
                <Slider {...fields.verify_ssl_certs} />
                <span className="slider-option slider-option--on">ON</span>
              </div>
            </div>

            <div className="input-field__wrapper">
              <label
                className="input-field__label"
                htmlFor={fields.enable_start_tls.name}
              >
                Enable STARTTLS?
              </label>
              <div className="slide-wrapper">
                <span className="slider-option slider-option--off">OFF</span>
                <Slider {...fields.enable_start_tls} />
                <span className="slider-option slider-option--on">ON</span>
              </div>
            </div>
          </div>
        </div>

        <div className={`${baseClass}__details`}>
          <p><strong>Domain</strong> - If you need to specify a HELO domain, you can do it here <em className="hint hint--brand">(Default: <strong>Blank</strong>)</em></p>
          <p><strong>Verify SSL Certs</strong> - Turn this off (not recommended) if you use a self-signed certificate <em className="hint hint--brand">(Default: <strong>On</strong>)</em></p>
          <p><strong>Enable STARTTLS</strong> - Detects if STARTTLS is enabled in your SMTP server and starts to use it. <em className="hint hint--brand">(Default: <strong>On</strong>)</em></p>
        </div>
      </div>
    );
  }

  render () {
    const { fields, handleSubmit, smtpConfigured } = this.props;
    const { onToggleAdvancedOptions, renderAdvancedOptions } = this;
    const { showAdvancedOptions } = this.state;

    return (
      <form className={baseClass}>
        <div className={`${baseClass}__section`}>
          <h2>Organization Info</h2>
          <div className={`${baseClass}__inputs`}>
            <InputField
              {...fields.org_name}
              label="Organization Name"
            />
            <InputField
              {...fields.org_logo_url}
              label="Organization Avatar"
            />
          </div>
          <div className={`${baseClass}__details ${baseClass}__avatar-preview`}>
            <img
              alt="Avatar preview"
              src={fields.org_logo_url.value}
            />
            <p>Avatar Preview</p>
          </div>
        </div>
        <div className={`${baseClass}__section`}>
          <h2>Kolide Web Address</h2>
          <div className={`${baseClass}__inputs`}>
            <InputField
              {...fields.kolide_server_url}
              label="Kolide App URL"
              hint={[`Include base path only (eg. no ${<code>/v1</code>})`]}
            />
          </div>
          <div className={`${baseClass}__details`}>
            <p>What base URL should <strong>osqueryd</strong> clients user to connect and register with <strong>Kolide</strong>?</p>
            <p className={`${baseClass}__note`}><strong>Note:</strong>Please ensure the URL you choose is accessible to all endpoints that need to communicate with Kolide, otherwise they will not be able to correctly register.</p>
            <Button text="SEND TEST" variant="inverse" />
          </div>
        </div>
        <div className={`${baseClass}__section`}>
          <h2>SMTP Options <small className={`smtp-options smtp-options--${smtpConfigured ? 'configured' : 'notconfigured'}`}>STATUS: <em>{smtpConfigured ? 'CONFIGURED' : 'NOT CONFIGURED'}</em></small></h2>
          <div className={`${baseClass}__inputs`}>
            <InputField
              {...fields.sender_address}
              label="Sender Address"
            />
          </div>
          <div className={`${baseClass}__details`}>
            <p>The address email recipients will see all messages that are sent from the <strong>Kolide</strong> application.</p>
          </div>
          <div className={`${baseClass}__inputs ${baseClass}__inputs--smtp`}>
            <InputField
              {...fields.server}
              label="SMTP Server"
            />
            <InputField
              {...fields.port}
              label="&nbsp;"
            />
            <Checkbox
              {...fields.enable_ssl_tls}
            >
              User SSL/TLS to connect (recommended)
            </Checkbox>
          </div>
          <div className={`${baseClass}__details`}>
            <p>The hostname / IP address and corresponding port of your organization&apos;s SMTP server.</p>
          </div>
          <div className={`${baseClass}__inputs`}>
            <div className="input-field__wrapper">
              <label
                className="input-field__label"
                htmlFor={fields.authentication_type.name}
              >
                Authentication Type
              </label>
              <Dropdown
                {...fields.authentication_type}
                options={authTypeOptions}
              />
            </div>
            <div className={`${baseClass}__smtp-section`}>
              <InputField
                {...fields.user_name}
                label="SMTP Username"
              />
              <InputField
                {...fields.password}
                label="SMTP Password"
              />
              <div className="input-field__wrapper">
                <label
                  className="input-field__label"
                  htmlFor={fields.auth_method.name}
                >
                  Auth Method
                </label>
                <Dropdown
                  {...fields.auth_method}
                  options={authMethodOptions}
                  placeholder=""
                />
              </div>
            </div>
          </div>
          <div className={`${baseClass}__details`}>
            <p>If your mail server requires authentication, you need to specify the authentication type here.</p>
            <p><strong>No Authentication</strong> - Select this if your SMTP is open.</p>
            <p><strong>Username & Password</strong> - Select this if your SMTP server requires username and password before use.</p>
          </div>
        </div>
        <div className={`${baseClass}__section`}>
          <h2><a href="#advancedOptions" onClick={onToggleAdvancedOptions} className={`${baseClass}__show-options`}><Header showAdvancedOptions={showAdvancedOptions} /></a></h2>

          {renderAdvancedOptions()}

        </div>
        <Button
          onClick={handleSubmit}
          text="UPDATE SETTINGS"
          variant="brand"
        />
      </form>
    );
  }
}

export default Form(AppConfigForm, {
  fields: formFields,
  validate,
});

