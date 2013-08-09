define(['jquery'], function($) {
	var agoltools = function() {
		// EMPTY
	};

	agoltools.toggler = function() {
		$('[agoltools-toggle]').bind('click', function() {
			$("#" + $(this).attr('agoltools-toggle')).toggleClass('hide');
		});
	};

	agoltools.expandAll = function() {
		$('[agoltools-toggle]').each(function() {
			$("#" + $(this).attr('agoltools-toggle')).removeClass('hide');
		});
	};

	agoltools.collapseAll = function() {
		$('[agoltools-toggle]').each(function() {
			$("#" + $(this).attr('agoltools-toggle')).addClass('hide');
		});
	};

	return agoltools;
});