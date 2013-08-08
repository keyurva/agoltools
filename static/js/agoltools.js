define(['jquery'], function($) {
	var agoltools = function() {
		// EMPTY
	};

	agoltools.toggler = function() {
		$('[agoltools-toggle]').bind('click', function() {
			$("#" + $(this).attr('agoltools-toggle')).toggleClass('hide');
		});
	};

	return agoltools;
});