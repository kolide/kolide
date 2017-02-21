import { flatMap, kebabCase, omit, pick, size } from 'lodash';
import md5 from 'js-md5';

const ORG_INFO_ATTRS = ['org_name', 'org_logo_url'];
const ADMIN_ATTRS = ['email', 'name', 'password', 'password_confirmation', 'username'];

export const addGravatarUrlToResource = (resource) => {
  const { email } = resource;

  const emailHash = md5(email.toLowerCase());
  const gravatarURL = `https://www.gravatar.com/avatar/${emailHash}?d=blank&size=200`;

  return {
    ...resource,
    gravatarURL,
  };
};

const labelSlug = (label) => {
  const { display_text: displayText } = label;

  if (!displayText) return undefined;

  const lowerDisplayText = displayText.toLowerCase();

  return kebabCase(lowerDisplayText);
};

const filterTarget = (targetType) => {
  return (target) => {
    return target.target_type === targetType ? [target.id] : [];
  };
};

export const formatConfigDataForServer = (config) => {
  const orgInfoAttrs = pick(config, ['org_logo_url', 'org_name']);
  const serverSettingsAttrs = pick(config, ['kolide_server_url', 'osquery_enroll_secret']);
  const smtpSettingsAttrs = pick(config, [
    'authentication_method', 'authentication_type', 'domain', 'email_enabled', 'enable_ssl_tls',
    'enable_start_tls', 'password', 'port', 'sender_address', 'server', 'user_name', 'verify_ssl_certs',
  ]);

  const orgInfo = size(orgInfoAttrs) && { org_info: orgInfoAttrs };
  const serverSettings = size(serverSettingsAttrs) && { server_settings: serverSettingsAttrs };
  const smtpSettings = size(smtpSettingsAttrs) && { smtp_settings: smtpSettingsAttrs };

  return {
    ...orgInfo,
    ...serverSettings,
    ...smtpSettings,
  };
};

export const formatSelectedTargetsForApi = (selectedTargets, appendID = false) => {
  const targets = selectedTargets || [];
  const hosts = flatMap(targets, filterTarget('hosts'));
  const labels = flatMap(targets, filterTarget('labels'));

  if (appendID) {
    return { host_ids: hosts, label_ids: labels };
  }

  return { hosts, labels };
};

const parseLicense = (license) => {
  const allowedHosts = license.allowed_hosts === 0 ? 'Unlimited' : license.allowed_hosts;

  return { ...license, allowed_hosts: allowedHosts };
};

export const formatScheduledQueryForServer = (scheduledQuery) => {
  const {
    interval,
    logging_type: loggingType,
    pack_id: packID,
    platform,
    query_id: queryID,
    shard,
  } = scheduledQuery;
  const result = omit(scheduledQuery, ['logging_type']);

  if (platform === 'all') {
    result.platform = '';
  }

  if (interval) {
    result.interval = Number(interval);
  }

  if (loggingType) {
    result.removed = loggingType === 'differential';
    result.snapshot = loggingType === 'snapshot';
  }

  if (packID) {
    result.pack_id = Number(packID);
  }

  if (queryID) {
    result.query_id = Number(queryID);
  }

  if (shard) {
    result.shard = Number(shard);
  }

  return result;
};

export const formatScheduledQueryForClient = (scheduledQuery) => {
  if (scheduledQuery.platform === '') {
    scheduledQuery.platform = 'all';
  }

  if (scheduledQuery.snapshot) {
    scheduledQuery.logging_type = 'snapshot';
  } else if (scheduledQuery.removed) {
    scheduledQuery.logging_type = 'differential';
  } else {
    scheduledQuery.logging_type = 'differential_ignore_removals';
  }

  return scheduledQuery;
};

const setupData = (formData) => {
  const orgInfo = pick(formData, ORG_INFO_ATTRS);
  const adminInfo = pick(formData, ADMIN_ATTRS);

  return {
    kolide_server_url: formData.kolide_server_url,
    org_info: {
      ...orgInfo,
    },
    admin: {
      admin: true,
      ...adminInfo,
    },
  };
};

export default {
  addGravatarUrlToResource,
  formatConfigDataForServer,
  formatScheduledQueryForClient,
  formatScheduledQueryForServer,
  formatSelectedTargetsForApi,
  labelSlug,
  parseLicense,
  setupData,
};
