(function(){
    $('.-about').parallax({imageSrc: './web/asset/img/about-x.jpg', zIndex: -1});

    $('.root').fullpage({
        easingcss3: 'cubic-bezier(0.770, 0.000, 0.175, 1.000)',
        loopHorizontal: false,
        continuousVertical: false,
        verticalCentered: false,
        resize : false,
        sectionSelector: '.full-section',
        scrollBar: true,
        navigation: true,
        afterResize: function(){
            $(window).trigger('resize.px.parallax')
        }
    });

    var $navLinks = $('.main-nav a');
    var sections = $('[data-anchor]').map(function(){
        return $(this).attr('data-anchor');
    }).toArray()

    $(window).on('hashchange', function(){
        var hash = window.location.hash;
        $navLinks.removeClass('-active');

        var index = sections.indexOf(hash.slice(1));
        if (index !== -1) {
            $navLinks.eq(index).addClass('-active');
        } else if (hash == '#about') {
            $navLinks.eq(-2).addClass('-active');
        }
    }).trigger('hashchange')
}())
