import React from 'react';
import expect from 'expect';
import { mount } from 'enzyme';

import { ConfigOptionsPage } from 'pages/config/ConfigOptionsPage/ConfigOptionsPage';
import { configOptionStub } from 'test/stubs';
import { fillInFormInput, itBehavesLikeAFormDropdownElement } from 'test/helpers';

describe('ConfigOptionsPage - component', () => {
  const blankConfigOption = { name: '', value: '' };
  const props = { configOptions: [] };

  it('renders', () => {
    const page = mount(<ConfigOptionsPage {...props} />);

    expect(page.length).toEqual(1);
  });

  it('renders reset and save buttons', () => {
    const page = mount(<ConfigOptionsPage {...props} />);
    const buttons = page.find('Button');
    const resetButton = buttons.find('.config-options-page__reset-btn');
    const saveButton = buttons.find('.config-options-page__save-btn');

    expect(resetButton.length).toEqual(1);
    expect(saveButton.length).toEqual(1);
  });

  describe('removing a config option', () => {
    it('sets the option value to null in state', () => {
      const page = mount(<ConfigOptionsPage configOptions={[configOptionStub]} />);
      const removeBtn = page.find('ConfigOptionForm').find('Button').first();

      expect(page.state('configOptions')).toEqual([configOptionStub]);

      removeBtn.simulate('click');

      expect(page.state('configOptions')).toEqual([{
        ...configOptionStub,
        value: null,
      }]);
    });
  });

  describe('adding a config option', () => {
    it('adds a blank option to state', () => {
      const page = mount(<ConfigOptionsPage configOptions={[configOptionStub]} />);
      const addBtn = page.find('Button').last();

      expect(page.state('configOptions')).toEqual([configOptionStub]);

      addBtn.simulate('click');

      expect(page.state('configOptions')).toEqual([
        configOptionStub,
        blankConfigOption,
      ]);
    });

    it('only allows one blank config option', () => {
      const page = mount(<ConfigOptionsPage configOptions={[configOptionStub]} />);
      const addBtn = page.find('Button').last();

      expect(page.state('configOptions')).toEqual([configOptionStub]);

      addBtn.simulate('click');
      addBtn.simulate('click');

      expect(page.state('configOptions')).toEqual([
        configOptionStub,
        blankConfigOption,
      ]);
    });
  });

  describe('updating a config option', () => {
    it('updates the config option in state', () => {
      const page = mount(<ConfigOptionsPage configOptions={[configOptionStub]} />);
      const configOptionInput = page.find('ConfigOptionForm').find('InputField');

      fillInFormInput(configOptionInput.find('input'), 'updated value');

      expect(page.state('configOptions')).toEqual([
        { ...configOptionStub, value: 'updated value' },
      ]);
    });

    it('updates the correct config option when the config option name is changed', () => {
      const configOptions = [configOptionStub, { id: 99, name: 'test', value: null }];
      const page = mount(<ConfigOptionsPage configOptions={configOptions} />);
      const form = page.find('ConfigOptionForm');
      const dropdownField = form.find('Select');
      const dropdownFieldInput = form.find('.Select-input').find('input');

      fillInFormInput(dropdownFieldInput, 'test');
      dropdownField.find('.Select-option').first().simulate('mousedown');

      expect(page.state('configOptions')).toEqual([
        configOptionStub,
        { id: 99, name: 'test', value: configOptionStub.value },
      ]);

      const configOptionInput = page.find('ConfigOptionForm').find('InputField').last();
      fillInFormInput(configOptionInput.find('input'), 'updated value');

      expect(page.state('configOptions')).toEqual([
        configOptionStub,
        { id: 99, name: 'test', value: 'updated value' },
      ]);
    });
  });
});
