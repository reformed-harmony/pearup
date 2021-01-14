$(function () {

  $('.users').each(function () {
    var $input = $(this);
    $input.autoComplete({
      minLength: 1,
      resolverSettings: {
        url: '/api/users'
      },
      formatResult: (item) => {
        return {
          value: item.id,
          text: item.text,
          html: [
            $('<img>').attr('src', '/media/' + item.picture).css({
              height: 18,
              marginRight: 12
            }),
            item.text
          ]
        };
      }
    });
    $input.on('autocomplete.select', function (_, i) {
      $input.siblings(
        'input[name=' + $input.data('control') + ']'
      ).val(i.value);
    })
  });

  // Forms with .pearup-confirm should display a confirmation dialog before
  // continuing with the specified action
  $('.pearup-confirm').each(function () {
    $(this).submit(function (e) {
      return window.confirm("Are you sure you wish to complete this action?");
    });
  });

});
