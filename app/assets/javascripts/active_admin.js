//= require active_admin/base
//= require rejection_modal
//= require select2-full

let locationSelect;

let select2Lang = {
  inputTooShort: function(args) {
    // args.minimum is the minimum required length
    // args.input is the user-typed text
    return  locationSelect.data('input-short');
  },
  errorLoading: function() {
    return locationSelect.data('error-loading');
  },
  loadingMore: function() {
    return locationSelect.data('loading-more');
  },
  noResults: function() {
    return locationSelect.data('no-results');
  },
  searching: function() {
    return locationSelect.data('on-searching');
  }
  }

$(document).ready(function () {
  locationSelect = $(".location-selector")
  locationSelect.select2({
    ajax: {
      url: "/api/v1/cities",
      delay: 500,
      dataType: 'json',
      data: function (params) {
        return {
          [locationSelect.data('searching-param')]: params.term // search term
        };
      },
      processResults: function (data) {
        return {
          results: mapCitiestoSelect2Array(data)
        };
      },
      cache: true
    },
    minimumInputLength: 3,
    allowClear: true,
    language:select2Lang
  });
});


function mapCitiestoSelect2Array(data){
  select2Array = []
  for (var i = 0; i < data.length; ++i) {
    select2Array.push({id: data[i].id, text: data[i][locationSelect.data('searching-param')]})
  }
  return select2Array;
}
