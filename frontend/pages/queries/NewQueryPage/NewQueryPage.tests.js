import expect, { restoreSpies } from 'expect';
import { mount } from 'enzyme';

import { defaultSelectedOsqueryTable } from '../../../redux/nodes/app/actions';
import helpers from '../../../test/helpers';
import NewQueryPage from './NewQueryPage';

const { connectedComponent, createAceSpy, reduxMockStore } = helpers;

describe('NewQueryPage - component', () => {
  beforeEach(createAceSpy);
  afterEach(restoreSpies);

  const mockStore = reduxMockStore({
    app: {
      selectedOsqueryTable: defaultSelectedOsqueryTable,
    },
    entities: {
      targets: {},
    },
  });

  it('renders the NewQuery component', () => {
    const page = mount(connectedComponent(NewQueryPage, { mockStore }));

    expect(page.find('NewQuery').length).toEqual(1);
  });

  it('renders the QuerySidePanel component', () => {
    const page = mount(connectedComponent(NewQueryPage, { mockStore }));

    expect(page.find('QuerySidePanel').length).toEqual(1);
  });
});
