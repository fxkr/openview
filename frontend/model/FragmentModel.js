export default class FragmentModel {
  constructor(urlModel) {
    this.urlModel = urlModel;
  }

  getFragments() {
    let fragmentNames = this.urlModel.getDirectory().split('/');
    if (fragmentNames[0] === '') {
      fragmentNames = fragmentNames.slice(1);
    }
    return fragmentNames.map((name, i) => ({
      name,
      url: this.urlModel.getDirectoryURL(fragmentNames.slice(0, i + 1).join('/')),
    }));
  }

  getHomeFragment() {
    return {
      name: 'Home',
      url: this.urlModel.getDirectoryURL(''),
    };
  }
}
